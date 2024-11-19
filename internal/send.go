package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Body string

	ChatID       string
	ChatThreadID string

	ConvertMarkdown bool
}

type SentMessage struct {
	MessageID int64
}

type tgSendResponse struct {
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}

type tgSendMessage struct {
	Text string `json:"text"`

	ChatID       string `json:"chat_id"`
	ChatThreadID string `json:"message_thread_id,omitempty"`

	Type   string `json:"type"`
	Markup string `json:"parse_mode"`
}

func (m *Messenger) Send(ctx context.Context, msg Message) (*SentMessage, error) {
	tgMsg := m.createTgMessage(msg)

	body, err := json.Marshal(tgMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message to json: %w", err)
	}

	uri := fmt.Sprintf("https://%s/bot%s/sendMessage", m.host, m.botToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	var res tgSendResponse

	err = m.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}

	return &SentMessage{
		MessageID: res.Result.MessageID,
	}, nil
}

func (m *Messenger) createTgMessage(msg Message) *tgSendMessage {
	tgMsg := &tgSendMessage{
		Text:         msg.Body,
		ChatID:       msg.ChatID,
		ChatThreadID: msg.ChatThreadID,
	}
	if !msg.ConvertMarkdown {
		return tgMsg
	}

	tgMsg.Type = "text"
	tgMsg.Markup = "HTML"

	tgMsg.Text = m.markdownConverter.Convert(msg.Body)

	return tgMsg
}
