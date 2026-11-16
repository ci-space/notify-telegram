package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cli "github.com/artarts36/singlecli"
	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendInlineKeyboard(t *testing.T) {
	tests := []struct {
		name   string
		mode   mode
		input  string
		method string
		body   string
	}{
		{
			name: "create preserves rows and button order",
			mode: modeCreate,
			input: `- - text: View PR
    url: https://github.com/org/repo/pull/123
  - text: CI Logs
    url: https://github.com/org/repo/actions
- - text: Open release
    url: https://github.com/org/repo/releases
`,
			method: "sendMessage",
			body: `{
				"chat_id":"chat","message_thread_id":"7","text":"hello","type":"","parse_mode":"",
				"reply_markup":{"inline_keyboard":[
					[{"text":"View PR","url":"https://github.com/org/repo/pull/123"},
					 {"text":"CI Logs","url":"https://github.com/org/repo/actions"}],
					[{"text":"Open release","url":"https://github.com/org/repo/releases"}]
				]}
			}`,
		},
		{
			name:   "update with one button",
			mode:   modeUpdate,
			input:  `- - {text: Open profile, url: 'tg://user?id=123'}`,
			method: "editMessageText",
			body: `{
				"chat_id":"chat","message_id":42,"text":"hello","type":"","parse_mode":"",
				"reply_markup":{"inline_keyboard":[[{"text":"Open profile","url":"tg://user?id=123"}]]}
			}`,
		},
		{
			name:   "create omits absent buttons",
			mode:   modeCreate,
			method: "sendMessage",
			body:   `{"chat_id":"chat","message_thread_id":"7","text":"hello","type":"","parse_mode":""}`,
		},
		{
			name:   "update omits blank buttons",
			mode:   modeUpdate,
			input:  " \n",
			method: "editMessageText",
			body:   `{"chat_id":"chat","message_id":42,"text":"hello","type":"","parse_mode":""}`,
		},
		{
			name:   "update omits an empty list",
			mode:   modeUpdate,
			input:  "[]",
			method: "editMessageText",
			body:   `{"chat_id":"chat","message_id":42,"text":"hello","type":"","parse_mode":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/bottoken/"+tt.method, r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				body, readErr := io.ReadAll(r.Body)
				if assert.NoError(t, readErr) {
					assert.JSONEq(t, tt.body, string(body))
				}
				_, writeErr := io.WriteString(w, `{"ok":true,"result":{"message_id":42}}`)
				assert.NoError(t, writeErr)
			}))
			t.Cleanup(server.Close)
			defaultClient := http.DefaultClient
			http.DefaultClient = server.Client()
			t.Cleanup(func() { http.DefaultClient = defaultClient })

			inputs := map[string]string{
				"INPUT_TOKEN":          "token",
				"INPUT_CHAT_ID":        "chat",
				"INPUT_CHAT_THREAD_ID": "7",
				"INPUT_MESSAGE":        "hello",
				"INPUT_HOST":           strings.TrimPrefix(server.URL, "https://"),
				"INPUT_MODE":           string(tt.mode),
			}
			if tt.mode == modeUpdate {
				inputs["INPUT_MESSAGE_ID"] = "42"
			}
			if tt.input != "" {
				inputs["INPUT_BUTTONS"] = tt.input
			}
			cfg, err := env.ParseAsWithOptions[config](env.Options{Environment: inputs})
			require.NoError(t, err)

			id, err := send(&cli.Context{Context: context.Background()}, cfg)
			require.NoError(t, err)
			require.EqualValues(t, 42, id)
		})
	}
}

func TestParseButtonsRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"invalid YAML", "[[", "yaml:"},
		{"multiple documents", "[]\n---\n[]", "single YAML document"},
		{"mapping instead of list", "buttons: []", "line 1: cannot unmarshal !!map"},
		{"null document", "null", "expected a list of rows"},
		{"flat list", `- {text: Open, url: https://example.com}`, "line 1: cannot unmarshal !!map"},
		{"empty row", "- []", "row 1:"},
		{"null row", "- null", "row 1:"},
		{"scalar button", "- - Open", "line 1: cannot unmarshal !!str"},
		{"null button", "- - null", "row 1, button 1: expected a mapping"},
		{"missing text", `- - {url: https://example.com}`, "row 1, button 1: text is required"},
		{"missing URL", `- - {text: Open}`, "row 1, button 1: url is required"},
		{"blank text", `- - {text: ' ', url: https://example.com}`, "row 1, button 1: text must be a non-empty string"},
		{"blank URL", `- - {text: Open, url: ''}`, "row 1, button 1: url must be a non-empty string"},
		{"non-string text", `- - {text: 123, url: https://example.com}`, "row 1, button 1: text must be a non-empty string"},
		{"boolean text", `- - {text: true, url: https://example.com}`, "row 1, button 1: text must be a non-empty string"},
		{"non-string URL", `- - {text: Open, url: [example.com]}`, "row 1, button 1: url must be a non-empty string"},
		{"unknown field", `- - text: Open
    url: https://example.com
    callback_data: data`, "line 3: field callback_data not found"},
		{"duplicate field", `- - text: Open
    text: Edit`, `line 2: mapping key "text" already defined at line 1`},
		{"relative URL", `- - {text: Open, url: example.com}`, "row 1, button 1: url must be a valid"},
		{"missing host", `- - {text: Open, url: 'https:///path'}`, "row 1, button 1: url must be a valid"},
		{"wrong scheme", `- - {text: Open, url: 'ftp://example.com'}`, "row 1, button 1: url must be a valid"},
		{"invalid URL escape",
			`- - {text: Open, url: 'https://example.com/?value=%zz'}`, "row 1, button 1: url must be a valid"},
		{"invalid button in second row", `[
			[{text: Valid, url: https://example.com}],
			[{text: Valid, url: https://example.com}, {text: Invalid, url: 'https://example.com/a b'}]
		]`, "row 2, button 2: url must be a valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyboard, err := parseButtons(tt.input)
			require.ErrorContains(t, err, tt.want)
			require.Nil(t, keyboard)
		})
	}
}
