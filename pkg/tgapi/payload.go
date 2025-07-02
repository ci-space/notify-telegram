package tgapi

type sendMessageRequest struct {
	Text string `json:"text"`

	ChatID       string `json:"chat_id"`
	ChatThreadID string `json:"message_thread_id,omitempty"`

	Type   string `json:"type"`
	Markup string `json:"parse_mode"`
}

type sendMessageResponse struct {
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}
