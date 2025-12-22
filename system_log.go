package unifi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"
)

var placeholderRegex = regexp.MustCompile(`\{(\w+)\}`)

// SystemLogEntry represents a v2 system log event from the UniFi Controller.
// API Path: /v2/api/site/{site}/system-log/all
type SystemLogEntry struct {
	Category        string                    `json:"category"`
	Event           string                    `json:"event"`
	ID              string                    `json:"id"`
	Key             string                    `json:"key"`
	MessageRaw      string                    `json:"message_raw"` // Used internally for Msg()
	Parameters      map[string]SystemLogParam `json:"parameters"`
	Severity        string                    `json:"severity"`
	ShowOnDashboard bool                      `json:"show_on_dashboard"`
	Status          string                    `json:"status"`
	Subcategory     string                    `json:"subcategory"`
	Target          string                    `json:"target"`
	Timestamp       int64                     `json:"timestamp"`
	TitleRaw        string                    `json:"title_raw"`
	Type            string                    `json:"type"`
	// Added by library
	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}

// SystemLogParam represents a parameter in a system log entry.
// Parameters contain contextual data like CLIENT, DEVICE, WLAN info.
type SystemLogParam struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	NotActionable       bool   `json:"not_actionable,omitempty"`
	Hostname            string `json:"hostname,omitempty"`
	IP                  string `json:"ip,omitempty"`
	Model               string `json:"model,omitempty"`
	ModelName           string `json:"model_name,omitempty"`
	Version             string `json:"version,omitempty"`
	DeviceFingerprintID int    `json:"device_fingerprint_id,omitempty"`
	FingerprintSource   int    `json:"fingerprint_source,omitempty"`
	NetworkPurpose      string `json:"network_purpose,omitempty"`
	Subnet              string `json:"subnet,omitempty"`
	VlanID              int    `json:"vlan_id,omitempty"`
}

// SystemLogRequest represents the request body for fetching system logs.
type SystemLogRequest struct {
	SearchText    string   `json:"searchText,omitempty"`
	Severities    []string `json:"severities,omitempty"`
	TimestampFrom int64    `json:"timestampFrom,omitempty"`
	TimestampTo   int64    `json:"timestampTo,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	Subcategories []string `json:"subcategories,omitempty"`
	Events        []string `json:"events,omitempty"`
	PageNumber    int      `json:"pageNumber"`
	PageSize      int      `json:"pageSize"`
}

// SystemLogResponse represents the response from the system log API.
type SystemLogResponse struct {
	Data              []*SystemLogEntry `json:"data"`
	PageNumber        int               `json:"page_number"`
	TotalElementCount int               `json:"total_element_count"`
	TotalPageCount    int               `json:"total_page_count"`
}

// DefaultSystemLogRequest returns a default request for fetching system logs
// from the last hour with all severities.
func DefaultSystemLogRequest(hours time.Duration) *SystemLogRequest {
	if hours < time.Hour {
		hours = time.Hour
	}

	now := time.Now()
	return &SystemLogRequest{
		Severities:    []string{"LOW", "MEDIUM", "HIGH", "VERY_HIGH"},
		TimestampFrom: now.Add(-hours).UnixMilli(),
		TimestampTo:   now.UnixMilli(),
		Categories:    []string{"MONITORING", "INTERNET", "POWER", "SECURITY", "SYSTEM"},
		PageNumber:    0,
		PageSize:      1000,
	}
}

// GetSystemLog returns system log events from multiple sites using the v2 API.
func (u *Unifi) GetSystemLog(sites []*Site, req *SystemLogRequest) ([]*SystemLogEntry, error) {
	data := make([]*SystemLogEntry, 0)

	for _, site := range sites {
		response, err := u.GetSiteSystemLog(site, req)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetSiteSystemLog retrieves system log events from a single site using the v2 API.
func (u *Unifi) GetSiteSystemLog(site *Site, req *SystemLogRequest) ([]*SystemLogEntry, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	if req == nil {
		req = DefaultSystemLogRequest(time.Hour)
	}

	u.DebugLog("Polling Controller for System Log (v2), site %s", site.SiteName)

	var allEntries []*SystemLogEntry
	currentPage := req.PageNumber

	// Paginate through all results
	for {
		reqCopy := *req
		reqCopy.PageNumber = currentPage

		path := fmt.Sprintf(APISystemLogPath, site.Name)

		// Marshal the request to JSON
		reqJSON, err := json.Marshal(reqCopy)
		if err != nil {
			return allEntries, fmt.Errorf("marshaling system log request: %w", err)
		}

		var response SystemLogResponse
		if err := u.GetData(path, &response, string(reqJSON)); err != nil {
			return allEntries, err
		}

		for _, entry := range response.Data {
			entry.SourceName = u.URL
			entry.SiteName = site.SiteName
			allEntries = append(allEntries, entry)
		}

		// Check if we've fetched all pages
		if currentPage >= response.TotalPageCount-1 || len(response.Data) == 0 {
			break
		}

		currentPage++

		// Safety limit to prevent infinite loops
		if currentPage > 100 {
			u.DebugLog("System log pagination limit reached (100 pages)")
			break
		}
	}

	// Sort by timestamp
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].Timestamp < allEntries[j].Timestamp
	})

	return allEntries, nil
}

// Datetime returns the timestamp as a time.Time for compatibility with Loki.
func (s *SystemLogEntry) Datetime() time.Time {
	return time.UnixMilli(s.Timestamp)
}

// Msg returns a formatted message from the system log entry.
// Replaces {PLACEHOLDER} tokens in MessageRaw with values from Parameters.
func (s *SystemLogEntry) Msg() string {
	if s.MessageRaw == "" {
		return s.TitleRaw
	}

	return placeholderRegex.ReplaceAllStringFunc(s.MessageRaw, func(match string) string {
		// Extract the placeholder name (remove { and })
		key := match[1 : len(match)-1]

		if param, ok := s.Parameters[key]; ok {
			if param.Name != "" {
				return param.Name
			}
			// Fallback to ID if Name is empty
			if param.ID != "" {
				return param.ID
			}
		}

		// Keep original placeholder if no replacement found
		return match
	})
}

// GetClientName extracts the client name from parameters if available.
func (s *SystemLogEntry) GetClientName() string {
	if client, ok := s.Parameters["CLIENT"]; ok {
		if client.Name != "" {
			return client.Name
		}
		return client.Hostname
	}
	return ""
}

// GetClientMAC extracts the client MAC address from parameters if available.
func (s *SystemLogEntry) GetClientMAC() string {
	if client, ok := s.Parameters["CLIENT"]; ok {
		return client.ID
	}
	return ""
}

// GetDeviceName extracts the device name from parameters if available.
func (s *SystemLogEntry) GetDeviceName() string {
	if device, ok := s.Parameters["DEVICE"]; ok {
		return device.Name
	}
	if device, ok := s.Parameters["DEVICE_TO"]; ok {
		return device.Name
	}
	return ""
}

// MarshalJSON customizes JSON output to include the formatted message
// and exclude message_raw (replaced by message).
func (s *SystemLogEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Category        string                    `json:"category"`
		Event           string                    `json:"event"`
		ID              string                    `json:"id"`
		Key             string                    `json:"key"`
		Message         string                    `json:"message"`
		Parameters      map[string]SystemLogParam `json:"parameters"`
		Severity        string                    `json:"severity"`
		ShowOnDashboard bool                      `json:"show_on_dashboard"`
		Status          string                    `json:"status"`
		Subcategory     string                    `json:"subcategory"`
		Target          string                    `json:"target"`
		Timestamp       int64                     `json:"timestamp"`
		TitleRaw        string                    `json:"title_raw"`
		Type            string                    `json:"type"`
		SiteName        string                    `json:"site_name,omitempty"`
		SourceName      string                    `json:"source_name,omitempty"`
	}{
		Category:        s.Category,
		Event:           s.Event,
		ID:              s.ID,
		Key:             s.Key,
		Message:         s.Msg(),
		Parameters:      s.Parameters,
		Severity:        s.Severity,
		ShowOnDashboard: s.ShowOnDashboard,
		Status:          s.Status,
		Subcategory:     s.Subcategory,
		Target:          s.Target,
		Timestamp:       s.Timestamp,
		TitleRaw:        s.TitleRaw,
		Type:            s.Type,
		SiteName:        s.SiteName,
		SourceName:      s.SourceName,
	})
}
