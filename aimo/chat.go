package aimo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) buildRequest(req *ChatCompletionRequest, stream bool) (*http.Request, error) {
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)

	fullReq := map[string]interface{}{
		"model":       c.model,
		"messages":    req.Messages,
		"stream":      stream,
		"max_tokens":  c.maxTokens,
		"temperature": c.temperature,
	}

	jsonData, err := json.Marshal(fullReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(c.ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	return httpReq, nil
}

func (c *Client) ChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	httpReq, err := c.buildRequest(req, false)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chatResp, nil
}

func (c *Client) ChatCompletionStream(req *ChatCompletionRequest) (<-chan StreamResponse, <-chan error) {
	responseChan := make(chan StreamResponse, 100)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		httpReq, err := c.buildRequest(req, true)
		if err != nil {
			errorChan <- err
			return
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			errorChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for {
			var streamResp StreamResponse
			if err := decoder.Decode(&streamResp); err != nil {
				if err == io.EOF {
					break
				}
				errorChan <- fmt.Errorf("failed to decode stream response: %w", err)
				return
			}
			responseChan <- streamResp
		}
	}()

	return responseChan, errorChan
}

func (c *Client) Chat(message string) (*ChatCompletionResponse, error) {
	req := &ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: message},
		},
	}
	return c.ChatCompletion(req)
}

func (c *Client) ChatStream(message string) (<-chan StreamResponse, <-chan error) {
	req := &ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: message},
		},
	}
	return c.ChatCompletionStream(req)
}

func (c *Client) ChatWithMessages(messages []string) (*ChatCompletionResponse, error) {
	msgList := make([]Message, len(messages))
	for i, msg := range messages {
		msgList[i] = Message{Role: "user", Content: msg}
	}

	req := &ChatCompletionRequest{
		Messages: msgList,
	}
	return c.ChatCompletion(req)
}

func (c *Client) ChatWithMessagesStream(messages []string) (<-chan StreamResponse, <-chan error) {
	msgList := make([]Message, len(messages))
	for i, msg := range messages {
		msgList[i] = Message{Role: "user", Content: msg}
	}

	req := &ChatCompletionRequest{
		Messages: msgList,
	}
	return c.ChatCompletionStream(req)
}
