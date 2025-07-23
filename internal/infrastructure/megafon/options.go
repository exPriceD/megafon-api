package megafon

import (
	"net/http"
	"time"
)

type ClientOption func(*Client)

func WithHTTPClient(h *http.Client) ClientOption {
	return func(c *Client) { c.http = h }
}

func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) { c.http.Timeout = d }
}

// WithRetries задаёт макс. число повторов (по-умолчанию 2 = 1 запрос + 1 повтор).
func WithRetries(n int) ClientOption {
	return func(c *Client) { c.retries = n }
}
