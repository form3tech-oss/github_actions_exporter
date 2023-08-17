package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type GrafanaClient struct {
	baseUrl url.URL
	opts    Opts
	client  *http.Client
}

func NewGrafanaClient(opts Opts) (*GrafanaClient, error) {
	url, err := url.Parse(opts.GrafanaBaseUrl)
	if err != nil {
		return nil, err
	}

	cli := http.DefaultClient

	return &GrafanaClient{
		baseUrl: *url,
		opts:    opts,
		client:  cli,
	}, nil
}

func (c *GrafanaClient) Request(ctx context.Context, method, reqPath string, query url.Values, result any) error {
	url := c.baseUrl
	url.Path = path.Join(url.Path, reqPath)
	url.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, method, url.String(), nil)
	if err != nil {
		return err
	}

	fmt.Println(url.String())

	if c.opts.GrafanaToken == "" {
		return errors.New("grafana token not set")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.opts.GrafanaToken))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var reason any
		_ = json.Unmarshal(b, &reason)
		return fmt.Errorf("unexpected error with status code %d, reason: %s", resp.StatusCode, reason)
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	return nil
}
