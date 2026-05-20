// Package wizebot provides a minimal HTTP client for the WizeBot ranking API.
package wizebot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL     = "https://wapi.wizebot.tv/api/ranking"
	httpTimeout = 15 * time.Second
)

// TopType identifies the ranking metric fetched from the API.
type TopType string

const (
	TopTypeUptime  TopType = "uptime"
	TopTypeMessage TopType = "message"
)

// RankingEntry is a single entry returned by the WizeBot ranking API.
type RankingEntry struct {
	UserName string  `json:"user_name"`
	UserUID  string  `json:"user_uid"`
	Value    float64 `json:"value"`
}

// rankingResponse is the raw API response envelope.
type rankingResponse struct {
	Success bool           `json:"success"`
	List    []RankingEntry `json:"list"`
}

// Client wraps HTTP calls to the WizeBot ranking API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// New returns a Client ready to query the WizeBot API with the provided key.
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// FetchTop retrieves the top entries for a given type and period from the API.
// topType is "uptime" or "message"; period is "week" or "month"; limit is 1-100.
func (c *Client) FetchTop(ctx context.Context, topType TopType, period string, limit int) ([]RankingEntry, error) {
	url := fmt.Sprintf("%s/%s/top/%s/%s/%d", baseURL, c.apiKey, topType, period, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for %s/%s: %w", topType, period, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s/%s: %w", topType, period, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WizeBot API returned HTTP %d for %s/%s", resp.StatusCode, topType, period)
	}

	var payload rankingResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding response for %s/%s: %w", topType, period, err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("WizeBot API reported failure for %s/%s", topType, period)
	}

	return payload.List, nil
}
