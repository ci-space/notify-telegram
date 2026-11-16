package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"unicode"

	cli "github.com/artarts36/singlecli"
	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"

	githuboutput "github.com/ci-space/github-output"
	"github.com/ci-space/notify-telegram/internal"
	"github.com/ci-space/notify-telegram/pkg/tgapi"
)

type config struct {
	Token           string `env:"INPUT_TOKEN,required"`
	ChatID          string `env:"INPUT_CHAT_ID,required"`
	ChatThreadID    string `env:"INPUT_CHAT_THREAD_ID"`
	Message         string `env:"INPUT_MESSAGE,required"`
	Buttons         string `env:"INPUT_BUTTONS"`
	Host            string `env:"INPUT_HOST,required"`
	Socks5Addr      string `env:"INPUT_SOCKS5_ADDR"`
	ConvertMarkdown bool   `env:"INPUT_CONVERT_MARKDOWN"`
	IssueTrackerURL string `env:"INPUT_ISSUE_TRACKER_URL"`
	Mode            mode   `env:"INPUT_MODE"`
	MessageID       int64  `env:"INPUT_MESSAGE_ID"`
}

type mode string

const (
	modeUnspecified = ""
	modeCreate      = "create"
	modeUpdate      = "update"
)

func main() {
	app := &cli.App{
		BuildInfo: &cli.BuildInfo{
			Name: "telegram-action",
		},
		Args:   []*cli.ArgDefinition{},
		Action: run,
	}

	app.RunWithGlobalArgs(context.Background())
}

const colorGreen = 2

func run(ctx *cli.Context) error {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Mode == modeUnspecified {
		cfg.Mode = modeCreate
	}

	if cfg.Mode == modeCreate && cfg.MessageID > 0 { //nolint:gocritic // not need
		return fmt.Errorf("message id must be empty, when mode = create")
	} else if cfg.Mode == modeUpdate && cfg.MessageID == 0 {
		return fmt.Errorf("message id must be filled, when mode = update")
	} else if !cfg.Mode.Valid() {
		return fmt.Errorf("mode %q unknown", cfg.Mode)
	}

	messageID, err := send(ctx, cfg)
	if err != nil {
		return err
	}

	ctx.Output.PrintColoredBlock(colorGreen, fmt.Sprintf("Message was sent. ID: %d", messageID))

	return githuboutput.WhenAvailable(func() error {
		return githuboutput.Write("message_id", fmt.Sprintf("%d", messageID))
	})
}

func send(ctx *cli.Context, cfg config) (int64, error) {
	buttons, err := parseButtons(cfg.Buttons)
	if err != nil {
		return 0, fmt.Errorf("invalid buttons: %w", err)
	}
	var keyboard *tgapi.InlineKeyboardMarkup
	if len(buttons) > 0 {
		keyboard = &tgapi.InlineKeyboardMarkup{InlineKeyboard: buttons}
	}

	client, err := tgapi.NewClient(tgapi.Config{
		BotToken:   cfg.Token,
		Host:       cfg.Host,
		Socks5Addr: cfg.Socks5Addr,
	})
	if err != nil {
		return 0, fmt.Errorf("create telegram client: %w", err)
	}

	switch cfg.Mode {
	case modeCreate:
		return create(ctx, client, cfg, keyboard)
	case modeUpdate:
		return update(ctx, client, cfg, keyboard)
	default:
		return 0, errors.New("unknown mode")
	}
}

func create(ctx *cli.Context, client *tgapi.Client, cfg config, keyboard *tgapi.InlineKeyboardMarkup) (int64, error) {
	msg := tgapi.SendingMessage{
		Body:            cfg.Message,
		ChatID:          cfg.ChatID,
		ChatThreadID:    cfg.ChatThreadID,
		ConvertMarkdown: cfg.ConvertMarkdown,
		ReplyMarkup:     keyboard,
	}

	if cfg.IssueTrackerURL != "" {
		tracker := internal.NewIssueTracker(cfg.IssueTrackerURL)
		msg.Body = tracker.InjectLinks(msg.Body)
	}

	res, err := client.SendMessage(ctx.Context, msg)
	if err != nil {
		return 0, fmt.Errorf("send message: %w", err)
	}

	return res.MessageID, nil
}

func update(ctx *cli.Context, client *tgapi.Client, cfg config, keyboard *tgapi.InlineKeyboardMarkup) (int64, error) {
	msg := tgapi.EditingMessageText{
		ID:              cfg.MessageID,
		ChatID:          cfg.ChatID,
		Body:            cfg.Message,
		ConvertMarkdown: cfg.ConvertMarkdown,
		ReplyMarkup:     keyboard,
	}

	if cfg.IssueTrackerURL != "" {
		tracker := internal.NewIssueTracker(cfg.IssueTrackerURL)
		msg.Body = tracker.InjectLinks(msg.Body)
	}

	res, err := client.EditMessageText(ctx.Context, msg)
	if err != nil {
		return 0, fmt.Errorf("edit message text: %w", err)
	}

	return res.MessageID, nil
}

func parseButtons(input string) ([][]tgapi.InlineKeyboardButton, error) {
	decoder := yaml.NewDecoder(strings.NewReader(input))
	var document yaml.Node
	err := decoder.Decode(&document)
	if errors.Is(err, io.EOF) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err = decoder.Decode(&yaml.Node{}); err == nil {
		return nil, errors.New("expected a single YAML document")
	}
	if !errors.Is(err, io.EOF) {
		return nil, err
	}

	rows := document.Content[0]
	if rows.Kind != yaml.SequenceNode {
		return nil, errors.New("expected a list of rows")
	}
	keyboard := make([][]tgapi.InlineKeyboardButton, len(rows.Content))
	for rowIndex, row := range rows.Content {
		if row.Kind != yaml.SequenceNode || len(row.Content) == 0 {
			return nil, fmt.Errorf("row %d: expected a non-empty list of buttons", rowIndex+1)
		}
		buttons := make([]tgapi.InlineKeyboardButton, len(row.Content))
		for buttonIndex, node := range row.Content {
			button, parseErr := parseButton(node)
			if parseErr != nil {
				return nil, fmt.Errorf("row %d, button %d: %w", rowIndex+1, buttonIndex+1, parseErr)
			}
			buttons[buttonIndex] = button
		}
		keyboard[rowIndex] = buttons
	}

	return keyboard, nil
}

func parseButton(node *yaml.Node) (tgapi.InlineKeyboardButton, error) {
	var button tgapi.InlineKeyboardButton
	if node.Kind != yaml.MappingNode {
		return button, errors.New("expected a mapping with text and url fields")
	}

	for i := 0; i < len(node.Content); i += 2 {
		key, value := node.Content[i], node.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.ShortTag() != "!!str" {
			return button, errors.New("expected a string field name")
		}

		var field *string
		switch key.Value {
		case "text":
			field = &button.Text
		case "url":
			field = &button.URL
		default:
			return button, fmt.Errorf("unsupported field %q", key.Value)
		}
		if *field != "" {
			return button, fmt.Errorf("duplicate field %q", key.Value)
		}
		if value.Kind != yaml.ScalarNode || value.ShortTag() != "!!str" || strings.TrimSpace(value.Value) == "" {
			return button, fmt.Errorf("%s must be a non-empty string", key.Value)
		}
		*field = value.Value
	}

	if button.Text == "" {
		return button, errors.New("text is required")
	}
	if button.URL == "" {
		return button, errors.New("url is required")
	}
	if !validButtonURL(button.URL) {
		return button, errors.New("url must be a valid http://, https://, or tg:// URL")
	}

	return button, nil
}

func validButtonURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" || strings.ContainsFunc(rawURL, unicode.IsSpace) {
		return false
	}
	if _, err = url.QueryUnescape(parsed.RawQuery); err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https" || parsed.Scheme == "tg"
}

func (m mode) Valid() bool {
	switch m {
	case modeCreate:
		return true
	case modeUpdate:
		return true
	default:
		return false
	}
}
