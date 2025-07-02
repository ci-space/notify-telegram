package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ci-space/notify-telegram/pkg/md2html"
	"net/http"
)

type EditingMessageText struct {
	ID              int64
	ChatID          string
	Body            string
	ConvertMarkdown bool
}

func (m *Client) EditMessageText(ctx context.Context, msg EditingMessageText) (*SentMessage, error) {
	tgMsg := m.createEditMessageTextRequest(msg)

	body, err := json.Marshal(tgMsg)
	if err != nil {
		return nil, fmt.Errorf("marshal message to json: %w", err)
	}

	uri := fmt.Sprintf("https://%s/bot%s/editMessageText", m.host, m.botToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	var res editMessageTextResponse

	err = m.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}

	return &SentMessage{
		MessageID: res.Result.MessageID,
	}, nil
}

func (m *Client) createEditMessageTextRequest(msg EditingMessageText) *editMessageTextRequest {
	tgMsg := &editMessageTextRequest{
		ID:     msg.ID,
		ChatID: msg.ChatID,
		Text:   msg.Body,
	}
	if !msg.ConvertMarkdown {
		return tgMsg
	}

	tgMsg.Type = "text"
	tgMsg.Markup = "HTML"

	tgMsg.Text = md2html.Convert(msg.Body)

	return tgMsg
}
