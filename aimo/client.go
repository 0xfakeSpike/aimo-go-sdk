package aimo

import (
	"context"
	"net/http"
	"time"
)

type Client struct {
	ctx         context.Context
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	model       string
	maxTokens   int
	temperature float64
}

func NewClient(ctx context.Context, endpoint, apiKey string) *Client {
	return &Client{
		ctx:         ctx,
		apiKey:      apiKey,
		baseURL:     endpoint,
		model:       "8W7X1tGnWh9CXwnPD7wgke31Gdcqmex4LapJvQ2afBUq:deepseek-chat-v3",
		maxTokens:   500,
		temperature: 0.7,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewClientWithConfig(ctx context.Context, endpoint, apiKey, model string, maxTokens int, temperature float64) *Client {
	return &Client{
		ctx:         ctx,
		apiKey:      apiKey,
		baseURL:     endpoint,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
