// Package wizebot provides a minimal HTTP client for the WizeBot ranking API.
package wizebot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// stringFloat64 is a float64 that unmarshals from either a JSON number or a
// JSON string (e.g. "1234" or "1234.5"). The WizeBot API returns value as a
// string-encoded number.
type stringFloat64 float64

func (s *stringFloat64) UnmarshalJSON(b []byte) error {
	raw := strings.Trim(string(b), `"`)
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("stringFloat64 : impossible de convertir %q : %w", raw, err)
	}
	*s = stringFloat64(v)
	return nil
}

// rawRankingEntry is used only for JSON decoding; its Value tolerates strings.
type rawRankingEntry struct {
	UserName string        `json:"user_name"`
	UserUID  string        `json:"user_uid"`
	Value    stringFloat64 `json:"value"`
}

// RankingEntry is a single entry returned by the WizeBot ranking API.
type RankingEntry struct {
	UserName string
	UserUID  string
	Value    float64
}

// rankingResponse is the raw API response envelope.
type rankingResponse struct {
	Success bool              `json:"success"`
	List    []rawRankingEntry `json:"list"`
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
		return nil, fmt.Errorf("construction de la requête %s/%s : %w", topType, period, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("appel API %s/%s : %w", topType, period, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API WizeBot a retourné HTTP %d pour %s/%s", resp.StatusCode, topType, period)
	}

	var payload rankingResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("décodage de la réponse %s/%s : %w", topType, period, err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("l'API WizeBot a signalé une erreur pour %s/%s", topType, period)
	}

	entries := make([]RankingEntry, len(payload.List))
	for i, r := range payload.List {
		entries[i] = RankingEntry{
			UserName: r.UserName,
			UserUID:  r.UserUID,
			Value:    float64(r.Value),
		}
	}
	return entries, nil
}
