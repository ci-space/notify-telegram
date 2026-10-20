package tgapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewClientWithoutSocks5Proxy(t *testing.T) {
	client, err := NewClient(Config{
		BotToken: "token",
		Host:     "api.telegram.org",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.httpClient != http.DefaultClient {
		t.Fatal("client without SOCKS5 proxy must use http.DefaultClient")
	}
}

func TestNewClientWithSocks5Proxy(t *testing.T) {
	client, err := NewClient(Config{
		BotToken:   "token",
		Host:       "api.telegram.org",
		Socks5Addr: "127.0.0.1:1080",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	transport, ok := client.httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.httpClient.Transport)
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.telegram.org", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	proxyURL, err := transport.Proxy(req)
	if err != nil {
		t.Fatalf("resolve proxy: %v", err)
	}
	if proxyURL == nil {
		t.Fatal("expected a proxy URL")
	}
	if proxyURL.String() != "socks5://127.0.0.1:1080" {
		t.Fatalf("unexpected proxy URL: %q", proxyURL.String())
	}
}

func TestNewClientWithInvalidSocks5Proxy(t *testing.T) {
	client, err := NewClient(Config{
		BotToken:   "token",
		Host:       "api.telegram.org",
		Socks5Addr: "http://proxy.example.com:1080",
	})

	if err == nil {
		t.Fatal("expected an error")
	}
	if client != nil {
		t.Fatal("expected a nil client")
	}
	if !strings.Contains(err.Error(), "invalid SOCKS5 proxy scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendMessageRedactsBotTokenFromRequestError(t *testing.T) {
	const token = "123456:secret-token"

	client := &Client{
		botToken: token,
		host:     "api.telegram.org",
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("request to %s failed", req.URL)
			}),
		},
	}

	_, err := client.SendMessage(context.Background(), SendingMessage{Body: "message", ChatID: "chat"})
	assertBotTokenRedacted(t, err, token)
}

func TestSendMessageRedactsBotTokenFromTelegramResponse(t *testing.T) {
	const token = "123456:secret-token"

	client := &Client{
		botToken: token,
		host:     "api.telegram.org",
		httpClient: &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				body := fmt.Sprintf(`{"ok":false,"description":"request to /bot%s/sendMessage failed"}`, token)

				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			}),
		},
	}

	_, err := client.SendMessage(context.Background(), SendingMessage{Body: "message", ChatID: "chat"})
	assertBotTokenRedacted(t, err, token)
}

func TestRedactBotTokenPreservesSafeError(t *testing.T) {
	originalErr := errors.New("request failed")

	if redactBotToken(originalErr, "token") != originalErr {
		t.Fatal("an error without the bot token must be preserved")
	}
}

func assertBotTokenRedacted(t *testing.T, err error, token string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error")
	}
	if strings.Contains(err.Error(), token) {
		t.Fatalf("bot token leaked in error: %s", err)
	}
	if !strings.Contains(err.Error(), redactedBotToken) {
		t.Fatalf("redaction marker is missing from error: %s", err)
	}
}

func TestParseSocks5ProxyURL(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		expected  string
		wantError string
	}{
		{
			name:     "address",
			addr:     "proxy.example.com:1080",
			expected: "socks5://proxy.example.com:1080",
		},
		{
			name:     "URL",
			addr:     "socks5://user:password@proxy.example.com:1080",
			expected: "socks5://user:password@proxy.example.com:1080",
		},
		{
			name:     "socks5h URL",
			addr:     "socks5h://proxy.example.com:1080",
			expected: "socks5h://proxy.example.com:1080",
		},
		{
			name:      "unsupported scheme",
			addr:      "http://proxy.example.com:1080",
			wantError: "invalid SOCKS5 proxy scheme",
		},
		{
			name:      "missing port",
			addr:      "proxy.example.com",
			wantError: "host and port are required",
		},
		{
			name:      "path",
			addr:      "socks5://proxy.example.com:1080/path",
			wantError: "path, query, and fragment are not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxyURL, err := parseSocks5ProxyURL(tt.addr)
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if proxyURL.String() != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, proxyURL.String())
			}
		})
	}
}
