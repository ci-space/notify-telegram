package tgapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const redactedBotToken = "[REDACTED]"

type Client struct {
	botToken   string
	host       string
	httpClient *http.Client
}

type Config struct {
	BotToken   string
	Host       string
	Socks5Addr string
}

func NewClient(cfg Config) (*Client, error) {
	client := &Client{
		botToken:   cfg.BotToken,
		host:       cfg.Host,
		httpClient: http.DefaultClient,
	}

	if strings.TrimSpace(cfg.Socks5Addr) == "" {
		return client, nil
	}

	proxyURL, err := parseSocks5ProxyURL(cfg.Socks5Addr)
	if err != nil {
		return nil, fmt.Errorf("configure SOCKS5 proxy: %w", err)
	}

	baseTransport := http.DefaultClient.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	defaultTransport, ok := baseTransport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("configure SOCKS5 proxy: unsupported HTTP transport %T", baseTransport)
	}

	transport := defaultTransport.Clone()
	transport.Proxy = http.ProxyURL(proxyURL)

	httpClient := *http.DefaultClient
	httpClient.Transport = transport
	client.httpClient = &httpClient

	return client, nil
}

func (m *Client) sendRequest(req *http.Request, out interface{}) error {
	req.Header.Add("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		resp.Body.Close()
	}()

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse body response: %w", err)
	}

	type errorResp struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	var errResp errorResp

	err = json.Unmarshal(bodyContent, &errResp)
	if err != nil {
		return err
	}

	if !errResp.Ok {
		return fmt.Errorf("telegram return no-ok response: %s", string(bodyContent))
	}

	err = json.Unmarshal(bodyContent, out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response json: %w", err)
	}

	return nil
}

func parseSocks5ProxyURL(addr string) (*url.URL, error) {
	addr = strings.TrimSpace(addr)
	if !strings.Contains(addr, "://") {
		addr = "socks5://" + addr
	}

	proxyURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid SOCKS5 proxy address: %w", err)
	}

	if proxyURL.Scheme != "socks5" && proxyURL.Scheme != "socks5h" {
		return nil, fmt.Errorf("invalid SOCKS5 proxy scheme %q", proxyURL.Scheme)
	}
	if proxyURL.Hostname() == "" || proxyURL.Port() == "" {
		return nil, fmt.Errorf("invalid SOCKS5 proxy address %q: host and port are required", addr)
	}
	if proxyURL.Path != "" || proxyURL.RawQuery != "" || proxyURL.Fragment != "" {
		return nil, fmt.Errorf("invalid SOCKS5 proxy address %q: path, query, and fragment are not supported", addr)
	}

	return proxyURL, nil
}

func redactBotToken(err error, token string) error {
	if err == nil || token == "" || !strings.Contains(err.Error(), token) {
		return err
	}

	return errors.New(strings.ReplaceAll(err.Error(), token, redactedBotToken))
}
