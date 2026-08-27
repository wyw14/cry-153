package gate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HTTP       *http.Client
	Endpoint   string
	Capacity   *Capacity
	MaxRetries int
	MaxBackoff time.Duration
}

func NewClient(endpoint string, capacity *Capacity) *Client {
	return &Client{HTTP: &http.Client{Timeout: 3 * time.Second}, Endpoint: endpoint, Capacity: capacity, MaxRetries: 3, MaxBackoff: 30 * time.Second}
}
func (c *Client) Close(ctx context.Context, intake, incident string) error {
	payload := []byte(fmt.Sprintf(`{"intake_id":"%s","incident_id":"%s","close":true}`, intake, incident))
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if err := c.Capacity.Acquire(ctx); err != nil {
			return err
		}
		err := c.send(ctx, payload)
		c.Capacity.Release()
		if err == nil {
			return nil
		}
		if attempt == c.MaxRetries {
			return err
		}
		if err := WaitBackoff(ctx, attempt, c.MaxBackoff); err != nil {
			return err
		}
	}
	return nil
}
func (c *Client) send(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("gateway status %d: %s", resp.StatusCode, string(body))
}
func WaitBackoff(ctx context.Context, attempt int, max time.Duration) error {
	d := time.Second * time.Duration(1<<attempt)
	if d > max || d < 0 {
		d = max
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
