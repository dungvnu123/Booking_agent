package rostering

import (
	"abclinic.local/backend/internal/platform/fault"
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client is for the existing Python OR-Tools solver, not an AI/LLM service.
type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: timeout}}
}

// Solve intentionally sends no request until the solver DTO contract is ported.
func (c *Client) Solve(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("rostering.Solve: %w", fault.ErrNotImplemented)
}
