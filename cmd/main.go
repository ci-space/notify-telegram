package main

import (
	"context"
	"errors"
	"fmt"

	cli "github.com/artarts36/singlecli"
	"github.com/caarlos0/env/v11"

	githuboutput "github.com/ci-space/github-output"
	"github.com/ci-space/notify-telegram/internal"
	"github.com/ci-space/notify-telegram/pkg/tgapi"
)

type config struct {
	Token           string `env:"INPUT_TOKEN,required"`
	ChatID          string `env:"INPUT_CHAT_ID,required"`
	ChatThreadID    string `env:"INPUT_CHAT_THREAD_ID"`
	Message         string `env:"INPUT_MESSAGE,required"`
	Host            string `env:"INPUT_HOST,required"`
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
	client := tgapi.NewClient(cfg.Token, cfg.Host)

	switch cfg.Mode {
	case modeCreate:
		return create(ctx, client, cfg)
	case modeUpdate:
		return update(ctx, client, cfg)
	default:
		return 0, errors.New("unknown mode")
	}
}

func create(ctx *cli.Context, client *tgapi.Client, cfg config) (int64, error) {
	msg := tgapi.SendingMessage{
		Body:            cfg.Message,
		ChatID:          cfg.ChatID,
		ChatThreadID:    cfg.ChatThreadID,
		ConvertMarkdown: cfg.ConvertMarkdown,
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

func update(ctx *cli.Context, client *tgapi.Client, cfg config) (int64, error) {
	msg := tgapi.EditingMessageText{
		ID:              cfg.MessageID,
		ChatID:          cfg.ChatID,
		Body:            cfg.Message,
		ConvertMarkdown: cfg.ConvertMarkdown,
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
