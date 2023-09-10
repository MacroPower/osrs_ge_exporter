package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string, client *http.Client) *Client {
	return &Client{
		baseURL: baseURL,
		client:  client,
	}
}

func (r *Client) Get(ctx context.Context, query string, params url.Values) ([]byte, int, error) {
	return r.doRequest(ctx, "GET", query, params, nil)
}

func (r *Client) doRequest(
	ctx context.Context, method, query string, params url.Values, buf io.Reader,
) ([]byte, int, error) {
	u, _ := url.Parse(r.baseURL)
	u.Path = path.Join(u.Path, query)
	if params != nil {
		u.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "https://github.com/MacroPower/osrs_ge_exporter")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed request: %w", err)
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return data, resp.StatusCode, nil
}
