package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Messenger struct {
	botToken string
	host     string

	markdownConverter *MarkdownToHTMLConverter
}

func NewMessenger(
	botToken string,
	host string,
	markdownRenderer *MarkdownToHTMLConverter,
) *Messenger {
	return &Messenger{botToken: botToken, host: host, markdownConverter: markdownRenderer}
}

func (m *Messenger) sendRequest(req *http.Request, out interface{}) error {
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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
