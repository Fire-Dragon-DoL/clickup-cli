package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/resolver"
)

const DefaultBaseURL = "https://api.clickup.com/api/v2"

type Client struct {
	apiKey     string
	baseURL    string
	spaceID    string
	httpClient *http.Client
}

func NewClient(apiKey, baseURL, spaceID string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		spaceID:    spaceID,
		httpClient: &http.Client{},
	}
}

func Do[Req any, Res any](c *Client, method, path string, body *Req) (Res, error) {
	var zero Res

	var reqBody *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return zero, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	var req *http.Request
	var err error
	if reqBody != nil {
		req, err = http.NewRequest(method, c.baseURL+path, reqBody)
	} else {
		req, err = http.NewRequest(method, c.baseURL+path, nil)
	}
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Err   string `json:"err"`
			ECODE string `json:"ECODE"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return zero, &Error{
			StatusCode: resp.StatusCode,
			Code:       errResp.ECODE,
			Message:    errResp.Err,
		}
	}

	var result Res
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

type Error struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("clickup api error (%d): %s [%s]", e.StatusCode, e.Message, e.Code)
}

// SearchTasks implements resolver.Searcher
func (c *Client) SearchTasks(query string) ([]resolver.SearchResult, error) {
	return nil, fmt.Errorf("search not implemented")
}

// SearchLists implements resolver.Searcher
func (c *Client) SearchLists(query string) ([]resolver.SearchResult, error) {
	return nil, fmt.Errorf("search not implemented")
}

// SearchFolders implements resolver.Searcher
func (c *Client) SearchFolders(query string) ([]resolver.SearchResult, error) {
	if c.spaceID == "" {
		return nil, fmt.Errorf("space ID is required to search folders")
	}

	folders, err := GetFolders(c, c.spaceID)
	if err != nil {
		return nil, err
	}

	var results []resolver.SearchResult
	for _, f := range folders {
		if f.Name == query {
			results = append(results, resolver.SearchResult{
				ID:   f.ID,
				Name: f.Name,
			})
		}
	}
	return results, nil
}

// SearchUsers implements resolver.Searcher
func (c *Client) SearchUsers(query string) ([]resolver.SearchResult, error) {
	return nil, fmt.Errorf("search not implemented")
}
