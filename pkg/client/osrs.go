package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type PriceClient struct {
	client *Client
}

func NewOSRSPriceClient() *PriceClient {
	return &PriceClient{
		client: NewClient("https://prices.runescape.wiki/api/v1/osrs", http.DefaultClient),
	}
}

func (c *PriceClient) GetLatest(ctx context.Context, params url.Values) (*DataLatest, error) {
	data := &DataLatest{}
	resp, _, err := c.client.Get(ctx, "latest", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest: %w", err)
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}

func (c *PriceClient) Get5m(ctx context.Context, params url.Values) (*DataAvg, error) {
	data := &DataAvg{}
	resp, _, err := c.client.Get(ctx, "5m", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get 5m: %w", err)
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}

func (c *PriceClient) Get1h(ctx context.Context, params url.Values) (*DataAvg, error) {
	data := &DataAvg{}
	resp, _, err := c.client.Get(ctx, "1h", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1h: %w", err)
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}

func (c *PriceClient) GetMapping(ctx context.Context, params url.Values) ([]ItemMapping, error) {
	data := []ItemMapping{}
	resp, _, err := c.client.Get(ctx, "mapping", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapping: %w", err)
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}
