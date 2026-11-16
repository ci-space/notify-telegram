package tgapi

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type sendMessageRequest struct {
	Text string `json:"text"`

	ChatID       string `json:"chat_id"`
	ChatThreadID string `json:"message_thread_id,omitempty"`

	Type   string `json:"type"`
	Markup string `json:"parse_mode"`

	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type sendMessageResponse struct {
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}

type editMessageTextRequest struct {
	ID     int64  `json:"message_id"`
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
	Type   string `json:"type"`
	Markup string `json:"parse_mode"`

	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type editMessageTextResponse struct {
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}
