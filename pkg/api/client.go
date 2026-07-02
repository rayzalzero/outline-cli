package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	tokenType  string
	csrfToken  string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	tokenType := detectTokenType(token)
	client := &Client{
		baseURL:   baseURL,
		token:     token,
		tokenType: tokenType,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				MaxIdleConnsPerHost: 10,
			},
		},
	}
	
	if tokenType == "jwt" {
		client.fetchCSRFToken()
	}
	
	return client
}

func detectTokenType(token string) string {
	if strings.HasPrefix(token, "ol_api_") {
		return "api_key"
	}
	if strings.Count(token, ".") == 2 {
		return "jwt"
	}
	return "jwt"
}

// Response is the standard Outline API response
type Response struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

// post makes a POST request to the Outline API
func (c *Client) post(endpoint string, payload interface{}) (*Response, error) {
	// Marshal payload
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal payload: %w", err)
		}
	} else {
		body = []byte("{}")
	}

	// Create request
	url := fmt.Sprintf("%s/api/%s", c.baseURL, endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Check API status
	if !apiResp.OK {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return &apiResp, nil
}

func (c *Client) fetchCSRFToken() {
	req, err := http.NewRequest("GET", c.baseURL, nil)
	if err != nil {
		return
	}
	
	req.Header.Set("Cookie", "accessToken="+c.token)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrf" {
			c.csrfToken = cookie.Value
			break
		}
	}
}
