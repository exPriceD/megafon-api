package megafon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"megafon-buisness-reports/internal/config"
	"megafon-buisness-reports/internal/interfaces"
	"net/http"
	"net/url"
	"time"
)

const userAgent = "megafon-reports/1.0"

type Client struct {
	base    *url.URL
	apiKey  string
	http    *http.Client
	log     interfaces.Logger
	retries int
}

func NewClient(cfg config.MegafonBuisnessConfig, lg interfaces.Logger, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse baseURL: %w", err)
	}
	c := &Client{
		base:   u,
		apiKey: cfg.APIKey,
		log:    lg,
		http: &http.Client{
			Timeout:   15 * time.Second,
			Transport: http.DefaultTransport,
		},
		retries: 2,
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// Do выполняет JSON-запрос с фиксированным числом повторов.
func (c *Client) Do(ctx context.Context, method, path string, q url.Values, in, out any) error {
	u := *c.base
	u.Path = c.base.ResolveReference(&url.URL{Path: path}).Path
	u.RawQuery = q.Encode()

	var reqBody []byte
	if in != nil {
		var err error
		reqBody, err = json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.retries; attempt++ {
		var body io.Reader
		if reqBody != nil {
			body = bytes.NewReader(reqBody)
		}

		req, _ := http.NewRequestWithContext(ctx, method, u.String(), body)
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("X-API-Key", c.apiKey)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.http.Do(req)
		if err != nil {
			c.log.Warn("transport error", "err", err, "attempt", attempt)
			lastErr = err
			continue
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			defer resp.Body.Close()

			if out != nil {
				if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
					return fmt.Errorf("decode: %w", err)
				}
			}
			return nil
		}

		b, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusTooManyRequests:
			return ErrRateLimited
		default:
			lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
			c.log.Warn("megafon retry", "status", resp.StatusCode, "attempt", attempt)
			time.Sleep(300 * time.Millisecond)
		}
	}
	return lastErr
}
