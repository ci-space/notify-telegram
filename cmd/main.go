package main

import (
	"context"
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

	res, err := send(ctx, cfg)
	if err != nil {
		return err
	}

	ctx.Output.PrintColoredBlock(colorGreen, fmt.Sprintf("SendingMessage was sent. ID: %d", res.MessageID))

	return githuboutput.WhenAvailable(func() error {
		return githuboutput.Write("message_id", fmt.Sprintf("%d", res.MessageID))
	})
}

func send(ctx *cli.Context, cfg config) (*tgapi.SentMessage, error) {
	client := tgapi.NewClient(cfg.Token, cfg.Host)

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
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return res, nil
}
