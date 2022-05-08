package mainlogic

import (
	"context"
	"net/http"
	"time"
)

type (
	client struct {
		baseURL string
		client  *http.Client
		timeout time.Duration
	}

	Fetch interface {
		FetchData(ctx context.Context, urlPath string, params map[string]string) (string, error)
	}
)

func New(baseURL string, httpClient *http.Client, timeout time.Duration) *client {
	return &client{
		baseURL: baseURL,
		client:  httpClient,
		timeout: timeout,
	}
}

func (c *client) FetchData(ctx context.Context, urlPath string, params map[string]string) (string, error) {
	return "", nil
}
