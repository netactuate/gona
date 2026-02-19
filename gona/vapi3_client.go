package gona

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	V3Version      = "0.1.0"
	V3BaseEndpoint = "https://vapi3.netactuate.com"
)

type V3Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
	userAgent  string
	debug      bool
}

type V3APIResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

type V3ListData struct {
	Data json.RawMessage  `json:"data"`
	Meta V3PaginationMeta `json:"meta"`
}

type V3PaginationMeta struct {
	Limit            int      `json:"limit"`
	Offset           int      `json:"offset"`
	Total            int      `json:"total"`
	SuggestedFilters []string `json:"suggestedFilters"`
}

type V3Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag,omitempty"`
}

type V3Capacity struct {
	AutoScaling bool `json:"autoscaling"`
	RequestedGB *int `json:"requestedGB"`
	TotalGB     int  `json:"totalGB"`
}

type V3Package struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewV3Client(apiKey, baseURL string) *V3Client {
	if baseURL == "" {
		baseURL = V3BaseEndpoint
	}

	transport := &http.Transport{
		TLSNextProto: make(
			map[string]func(string, *tls.Conn) http.RoundTripper,
		),
	}

	endpoint, _ := url.Parse(baseURL)

	return &V3Client{
		baseURL:   endpoint,
		apiKey:    apiKey,
		userAgent: "gona-v3/" + V3Version,
		debug:     os.Getenv("NA_API_DEBUG") != "",
		httpClient: &http.Client{
			Transport: transport,
		},
	}
}

func (c *V3Client) debugLog(format string, v ...any) {
	if !c.debug {
		return
	}
	log.Printf("[DEBUG][V3] "+format, v...)
}

func (c *V3Client) v3ApiKeyPath(path string) string {
	if strings.Contains(path, "?") {
		return path + "&key=" + c.apiKey
	}
	return path + "?key=" + c.apiKey
}

func (c *V3Client) doRequest(method, path string, body interface{}) (*V3APIResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		c.debugLog("%s request body for %s: %s", method, path, string(jsonData))
		bodyReader = bytes.NewBuffer(jsonData)
	}

	relPath, err := url.Parse(c.v3ApiKeyPath(path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse path %q: %w", path, err)
	}

	fullURL := c.baseURL.ResolveReference(relPath)

	req, err := http.NewRequest(method, fullURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.debugLog("%s %s", method, fullURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.debugLog("response status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode == 204 {
		return &V3APIResponse{Code: 204}, nil
	}

	if isSemanticNotFound(resp.StatusCode, string(respBody)) {
		return nil, &V3NotFoundError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var apiResp V3APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response (status %d): %s", resp.StatusCode, string(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, respBody, "", "  "); err == nil {
			return &apiResp, fmt.Errorf("API error on %s %s: HTTP %d\n%s",
				method, path, resp.StatusCode, pretty.String())
		}
		return &apiResp, fmt.Errorf("API error on %s %s: HTTP %d\n%s",
			method, path, resp.StatusCode, string(respBody))
	}

	if apiResp.Code != 0 && apiResp.Code != 200 {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, respBody, "", "  "); err == nil {
			return &apiResp, fmt.Errorf("API error on %s %s: code %d\n%s",
				method, path, apiResp.Code, pretty.String())
		}
		return &apiResp, fmt.Errorf("API error on %s %s: code %d\n%s",
			method, path, apiResp.Code, string(respBody))
	}

	return &apiResp, nil
}

func (c *V3Client) get(path string) (*V3APIResponse, error) {
	return c.doRequest("GET", path, nil)
}

func (c *V3Client) post(path string, body interface{}) (*V3APIResponse, error) {
	return c.doRequest("POST", path, body)
}

func (c *V3Client) patch(path string, body interface{}) (*V3APIResponse, error) {
	return c.doRequest("PATCH", path, body)
}

func (c *V3Client) put(path string, body interface{}) (*V3APIResponse, error) {
	return c.doRequest("PUT", path, body)
}

func (c *V3Client) del(path string) (*V3APIResponse, error) {
	return c.doRequest("DELETE", path, nil)
}

type V3NotFoundError struct {
	StatusCode int
	Body       string
}

func (e *V3NotFoundError) Error() string {
	return fmt.Sprintf("resource not found (HTTP %d): %s", e.StatusCode, e.Body)
}

func IsV3NotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*V3NotFoundError)
	return ok
}

func isSemanticNotFound(statusCode int, body string) bool {
	if statusCode == 404 || statusCode == 410 {
		return true
	}

	lower := strings.ToLower(body)
	notFoundMessages := []string{
		"does not exist",
		"not associated with your account",
		"not found",
	}

	for _, msg := range notFoundMessages {
		if strings.Contains(lower, msg) {
			return true
		}
	}

	return false
}