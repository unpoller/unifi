package unifi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// RemoteAPIBaseURL is the base URL for the UniFi remote API.
	RemoteAPIBaseURL = "https://api.ui.com"
	// RemoteAPIVersion is the API version for the remote API.
	RemoteAPIVersion = "v1"
)

// Console represents a UniFi console from the remote API.
type Console struct {
	ID            string `json:"id"`
	IPAddress     string `json:"ipAddress"`
	Type          string `json:"type"`
	Owner         bool   `json:"owner"`
	IsBlocked     bool   `json:"isBlocked"`
	ReportedState struct {
		Name     string `json:"name"`
		Hostname string `json:"hostname"`
		IP       string `json:"ip"`
		State    string `json:"state"`
		Mac      string `json:"mac"`
	} `json:"reportedState"`
	ConsoleName string `json:"-"` // Derived field: name from reportedState
}

// HostsResponse represents the response from /v1/hosts endpoint.
type HostsResponse struct {
	Data           []Console `json:"data"`
	HTTPStatusCode int      `json:"httpStatusCode"`
	TraceID        string   `json:"traceId"`
	NextToken      string   `json:"nextToken,omitempty"`
}

// RemoteSite represents a site from the remote API.
type RemoteSite struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SitesResponse represents the response from the sites endpoint.
type SitesResponse struct {
	Data           []RemoteSite `json:"data"`
	HTTPStatusCode int          `json:"httpStatusCode"`
	TraceID        string       `json:"traceId"`
}

// APIErrorResponse represents an error response from the remote API.
type APIErrorResponse struct {
	Code           string `json:"code"`
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	TraceID        string `json:"traceId"`
}

// RemoteAPIClient handles HTTP requests to the remote UniFi API.
// This is separate from the main Unifi client because it uses a different
// authentication method (X-API-Key header) and different base URL.
type RemoteAPIClient struct {
	apiKey   string
	baseURL  string
	client   *http.Client
	ErrorLog Logger // Optional, can be nil
	DebugLog Logger // Optional, can be nil
	Log      Logger // Optional, can be nil
}

// NewRemoteAPIClient creates a new remote API client.
// Logger functions are optional and can be nil (will default to discardLogs).
func NewRemoteAPIClient(apiKey string, errorLog, debugLog, log Logger) *RemoteAPIClient {
	if errorLog == nil {
		errorLog = discardLogs
	}
	if debugLog == nil {
		debugLog = discardLogs
	}
	if log == nil {
		log = discardLogs
	}

	return &RemoteAPIClient{
		apiKey:  apiKey,
		baseURL: RemoteAPIBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
		ErrorLog: errorLog,
		DebugLog: debugLog,
		Log:      log,
	}
}

// makeRequest makes an HTTP request to the remote API.
func (c *RemoteAPIClient) makeRequest(method, path string, queryParams map[string]string) ([]byte, error) {
	fullURL := c.baseURL + path

	if len(queryParams) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("parsing URL: %w", err)
		}

		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	c.DebugLog("Making %s request to: %s", method, fullURL)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse the error response for better error messages
		var apiErr APIErrorResponse
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			// Build a helpful error message based on status code
			errMsg := fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, apiErr.Message)
			
			// Add helpful suggestions for common error cases
			if resp.StatusCode == 403 || resp.StatusCode == 404 {
				suggestions := c.getErrorSuggestions(resp.StatusCode, apiErr.Message, path)
				if suggestions != "" {
					errMsg += "\n" + suggestions
				}
			}
			
			return nil, fmt.Errorf(errMsg)
		}
		
		// Fallback to generic error if we can't parse the response
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// getErrorSuggestions provides helpful suggestions for common API errors.
func (c *RemoteAPIClient) getErrorSuggestions(statusCode int, message, path string) string {
	var suggestions []string
	
	if statusCode == 403 {
		// Check if this is a sites endpoint (which might indicate NVR or firmware issues)
		if strings.Contains(path, "/sites") {
			suggestions = append(suggestions, 
				"  • This API key may be associated with a UniFi Protect (NVR) console, which does not support Network API endpoints.",
				"  • Ensure the API key was created for a UniFi Network console, not a UniFi Protect console.",
				"  • Verify the console firmware is compatible with the Network API (UDM/UDM-Pro/Cloud Gateway with Network application).",
				"  • Check that the console has the Network application installed and running.",
			)
		} else {
			suggestions = append(suggestions,
				"  • Verify the API key has the necessary permissions for this operation.",
				"  • Ensure the API key was created for the correct console type.",
				"  • Check that the console firmware version supports this API endpoint.",
			)
		}
	} else if statusCode == 404 {
		if strings.Contains(path, "/sites") {
			suggestions = append(suggestions,
				"  • The console may not have the Network application installed or enabled.",
				"  • Verify the console firmware version is compatible with the Network API.",
				"  • This endpoint may not be available for this console type (e.g., UniFi Protect/NVR consoles).",
			)
		} else {
			suggestions = append(suggestions,
				"  • The requested endpoint may not be available for this console type or firmware version.",
				"  • Verify the console firmware is up to date and supports this API endpoint.",
			)
		}
	}
	
	if len(suggestions) > 0 {
		return "Possible solutions:\n" + strings.Join(suggestions, "\n")
	}
	
	return ""
}

// DiscoverConsoles discovers all consoles available via the remote API.
// It handles pagination automatically and filters for console type only.
func (c *RemoteAPIClient) DiscoverConsoles() ([]Console, error) {
	// Start with first page
	queryParams := map[string]string{
		"pageSize": "10",
	}

	var allConsoles []Console
	nextToken := ""

	for {
		if nextToken != "" {
			queryParams["nextToken"] = nextToken
		} else {
			// Remove nextToken from params for first request
			delete(queryParams, "nextToken")
		}

		body, err := c.makeRequest("GET", "/v1/hosts", queryParams)
		if err != nil {
			return nil, fmt.Errorf("fetching consoles: %w", err)
		}

		var response HostsResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("parsing consoles response: %w", err)
		}

		// Filter for console type only
		for _, console := range response.Data {
			if console.Type == "console" && !console.IsBlocked {
				// Extract the console name from reportedState
				console.ConsoleName = console.ReportedState.Name
				if console.ConsoleName == "" {
					console.ConsoleName = console.ReportedState.Hostname
				}
				allConsoles = append(allConsoles, console)
			}
		}

		// Check if there's a nextToken to continue pagination
		if response.NextToken == "" {
			break
		}

		nextToken = response.NextToken
		c.DebugLog("Fetching next page of consoles with nextToken: %s", nextToken)
	}

	return allConsoles, nil
}

// DiscoverSites discovers all sites for a given console ID.
func (c *RemoteAPIClient) DiscoverSites(consoleID string) ([]RemoteSite, error) {
	path := fmt.Sprintf("/v1/connector/consoles/%s/proxy/network/integration/v1/sites", consoleID)

	queryParams := map[string]string{
		"offset": "0",
		"limit":  "100",
	}

	body, err := c.makeRequest("GET", path, queryParams)
	if err != nil {
		return nil, fmt.Errorf("fetching sites for console %s: %w", consoleID, err)
	}

	var response SitesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing sites response: %w", err)
	}

	return response.Data, nil
}
