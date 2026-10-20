package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ci-space/notify-telegram/pkg/md2html"
)

type SendingMessage struct {
	Body string

	ChatID       string
	ChatThreadID string

	ConvertMarkdown bool
}

type SentMessage struct {
	MessageID int64
}

func (m *Client) SendMessage(ctx context.Context, msg SendingMessage) (*SentMessage, error) {
	tgMsg := m.createSendMessageRequest(msg)

	body, err := json.Marshal(tgMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message to json: %w", err)
	}

	uri := fmt.Sprintf("https://%s/bot%s/sendMessage", m.host, m.botToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return nil, redactBotToken(fmt.Errorf("failed to make request: %w", err), m.botToken)
	}

	var res sendMessageResponse

	err = m.sendRequest(req, &res)
	if err != nil {
		return nil, redactBotToken(err, m.botToken)
	}

	return &SentMessage{
		MessageID: res.Result.MessageID,
	}, nil
}

func (m *Client) createSendMessageRequest(msg SendingMessage) *sendMessageRequest {
	tgMsg := &sendMessageRequest{
		Text:         msg.Body,
		ChatID:       msg.ChatID,
		ChatThreadID: msg.ChatThreadID,
	}
	if !msg.ConvertMarkdown {
		return tgMsg
	}

	tgMsg.Type = "text"
	tgMsg.Markup = "HTML"

	tgMsg.Text = md2html.Convert(msg.Body)

	return tgMsg
}
