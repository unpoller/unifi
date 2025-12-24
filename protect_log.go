package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var protectPlaceholderRegex = regexp.MustCompile(`\{(\w+)\}`)

// ProtectLogEntry represents a system log event from UniFi Protect.
// API Path: /proxy/protect/api/events/system-logs
type ProtectLogEntry struct {
	ID                string                 `json:"id"`
	ModelKey          string                 `json:"modelKey"`
	Type              string                 `json:"type"`
	Start             int64                  `json:"start"`
	End               int64                  `json:"end"`
	Score             int                    `json:"score"`
	SmartDetectTypes  []string               `json:"smartDetectTypes"`
	SmartDetectEvents []string               `json:"smartDetectEvents"`
	Camera            string                 `json:"camera"`
	Partition         *string                `json:"partition"`
	User              string                 `json:"user"`
	Metadata          *ProtectLogMetadata    `json:"metadata"`
	Thumbnail         string                 `json:"thumbnail"`
	Heatmap           string                 `json:"heatmap"`
	Timestamp         int64                  `json:"timestamp"`
	IsFavorite        bool                   `json:"isFavorite"`
	FavoriteObjectIDs []string               `json:"favoriteObjectIds"`
	Description       *ProtectLogDescription `json:"description"`
	Category          *string                `json:"category"`
	Aggregated        bool                   `json:"_aggregated"`
	AggregatedCount   int                    `json:"_count"`
	AggregatedEvents  []string               `json:"_aggregatedEvents"`
	// Added by library
	SourceName      string `json:"-"`
	ThumbnailBase64 string `json:"-"` // Base64-encoded thumbnail image (populated by caller)
}

// ProtectLogMetadata contains metadata for a Protect log entry.
type ProtectLogMetadata struct {
	ClientData         *ProtectClientData    `json:"clientData,omitempty"`
	UserAction         string                `json:"userAction,omitempty"`
	ClientPlatform     string                `json:"clientPlatform,omitempty"`
	RAMDescription     string                `json:"ramDescription,omitempty"`
	RAMClassifications []string              `json:"ramClassifications,omitempty"`
	ZonesStatus        map[string]ZoneStatus `json:"zonesStatus,omitempty"`
	DetectedAreas      []DetectedArea        `json:"detectedAreas,omitempty"`
	DetectedThumbnails []DetectedThumbnail   `json:"detectedThumbnails,omitempty"`
}

// ProtectClientData contains client-specific data in metadata.
type ProtectClientData struct {
	End     int64  `json:"end"`
	Start   int64  `json:"start"`
	EventID string `json:"eventId"`
}

// ZoneStatus represents the status of a detection zone.
type ZoneStatus struct {
	Level  int    `json:"level"`
	Status string `json:"status"`
}

// DetectedArea represents a detected area in the image.
type DetectedArea struct {
	AreaIndexes       []int  `json:"areaIndexes"`
	SmartDetectObject string `json:"smartDetectObject"`
}

// DetectedThumbnail represents a detected object thumbnail.
type DetectedThumbnail struct {
	Type          string                 `json:"type"`
	Coord         []int                  `json:"coord,omitempty"`
	ObjectID      string                 `json:"objectId,omitempty"`
	CroppedID     string                 `json:"croppedId,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Confidence    int                    `json:"confidence,omitempty"`
	ClockBestWall int64                  `json:"clockBestWall,omitempty"`
}

// ProtectLogDescription contains the description and message for a Protect log entry.
type ProtectLogDescription struct {
	EventMetadata *ProtectEventMetadata `json:"eventMetadata"`
	MessageRaw    string                `json:"messageRaw"`
	MessageKeys   []ProtectMessageKey   `json:"messageKeys"`
}

// ProtectEventMetadata contains event metadata.
type ProtectEventMetadata struct {
	EventType          string `json:"eventType"`
	Category           string `json:"category"`
	SubCategory        string `json:"subCategory"`
	Severity           string `json:"severity"`
	Title              string `json:"title"`
	DeviceEventClassID int    `json:"deviceEventClassId"`
}

// ProtectMessageKey represents a message placeholder key.
type ProtectMessageKey struct {
	Key    string                 `json:"key"`
	Text   string                 `json:"text"`
	Style  []string               `json:"style"`
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// ProtectLogRequest represents the query parameters for fetching Protect logs.
type ProtectLogRequest struct {
	Devices                []string
	Start                  int64
	End                    int64
	Limit                  int
	OrderDirection         string
	SmartDetectAggregation bool
	WithoutDescriptions    bool
}

// ProtectLogResponse represents the response from the Protect log API.
type ProtectLogResponse struct {
	Items []*ProtectLogEntry `json:"items"`
}

// DefaultProtectLogRequest returns a default request for fetching Protect logs
// from the specified duration (minimum 1 hour, defaults to 24 hours if 0).
func DefaultProtectLogRequest(duration time.Duration) *ProtectLogRequest {
	if duration < time.Hour {
		duration = 24 * time.Hour // Default to 24 hours for Protect logs
	}

	now := time.Now()

	return &ProtectLogRequest{
		Start:                  now.Add(-duration).UnixMilli(),
		End:                    now.UnixMilli(),
		Limit:                  1000, // Increased limit for larger time window
		OrderDirection:         "desc",
		SmartDetectAggregation: true,
		WithoutDescriptions:    false,
	}
}

// GetProtectLogs returns Protect system log events.
func (u *Unifi) GetProtectLogs(req *ProtectLogRequest) ([]*ProtectLogEntry, error) {
	if req == nil {
		req = DefaultProtectLogRequest(time.Hour)
	}

	u.DebugLog("Polling Controller for Protect Logs")

	// Build query parameters
	params := url.Values{}
	params.Set("start", strconv.FormatInt(req.Start, 10))
	params.Set("end", strconv.FormatInt(req.End, 10))
	params.Set("limit", strconv.Itoa(req.Limit))
	params.Set("orderDirection", req.OrderDirection)
	params.Set("smartDetectAggregation", strconv.FormatBool(req.SmartDetectAggregation))
	params.Set("withoutDescriptions", strconv.FormatBool(req.WithoutDescriptions))

	for _, device := range req.Devices {
		params.Add("devices", device)
	}

	path := APIProtectLogPath + "?" + params.Encode()

	var response ProtectLogResponse
	if err := u.GetData(path, &response); err != nil {
		return nil, fmt.Errorf("fetching protect logs: %w", err)
	}

	// Add source name to each entry
	for _, entry := range response.Items {
		entry.SourceName = u.URL
	}

	// Sort by timestamp (ascending order for consistency with other log types)
	sort.Slice(response.Items, func(i, j int) bool {
		return response.Items[i].Timestamp < response.Items[j].Timestamp
	})

	return response.Items, nil
}

// Datetime returns the timestamp as a time.Time for compatibility with Loki.
func (p *ProtectLogEntry) Datetime() time.Time {
	return time.UnixMilli(p.Timestamp)
}

// Msg returns a formatted message from the Protect log entry.
// Replaces {placeholder} tokens in MessageRaw with values from MessageKeys.
func (p *ProtectLogEntry) Msg() string {
	if p.Description == nil {
		return p.Type
	}

	if p.Description.MessageRaw == "" {
		if p.Description.EventMetadata != nil {
			return p.Description.EventMetadata.Title
		}

		return p.Type
	}

	// Build a map of key -> text for replacements
	keyMap := make(map[string]string)
	for _, mk := range p.Description.MessageKeys {
		keyMap[mk.Key] = mk.Text
	}

	return protectPlaceholderRegex.ReplaceAllStringFunc(p.Description.MessageRaw, func(match string) string {
		// Extract the placeholder name (remove { and })
		key := match[1 : len(match)-1]

		if text, ok := keyMap[key]; ok {
			return text
		}

		// Keep original placeholder if no replacement found
		return match
	})
}

// GetEventType returns the event type from description metadata.
func (p *ProtectLogEntry) GetEventType() string {
	if p.Description != nil && p.Description.EventMetadata != nil {
		return p.Description.EventMetadata.EventType
	}

	return p.Type
}

// GetCategory returns the category from description metadata.
func (p *ProtectLogEntry) GetCategory() string {
	if p.Description != nil && p.Description.EventMetadata != nil {
		return p.Description.EventMetadata.Category
	}

	if p.Category != nil {
		return *p.Category
	}

	return ""
}

// GetSubCategory returns the subcategory from description metadata.
func (p *ProtectLogEntry) GetSubCategory() string {
	if p.Description != nil && p.Description.EventMetadata != nil {
		return p.Description.EventMetadata.SubCategory
	}

	return ""
}

// GetSeverity returns the severity from description metadata.
func (p *ProtectLogEntry) GetSeverity() string {
	if p.Description != nil && p.Description.EventMetadata != nil {
		return p.Description.EventMetadata.Severity
	}

	return ""
}

// GetTitle returns the title from description metadata.
func (p *ProtectLogEntry) GetTitle() string {
	if p.Description != nil && p.Description.EventMetadata != nil {
		return p.Description.EventMetadata.Title
	}

	return ""
}

// GetCameraName extracts the camera name from message keys if available.
func (p *ProtectLogEntry) GetCameraName() string {
	if p.Description == nil {
		return ""
	}

	for _, mk := range p.Description.MessageKeys {
		if mk.Key == "deviceLink" || mk.Action == "viewDeviceDetails" {
			return mk.Text
		}
	}

	return ""
}

// GetUserName extracts the user name from message keys if available.
func (p *ProtectLogEntry) GetUserName() string {
	if p.Description == nil {
		return ""
	}

	for _, mk := range p.Description.MessageKeys {
		if mk.Key == "userLink" || mk.Action == "viewUsers" {
			return mk.Text
		}
	}

	return ""
}

// ThumbnailData holds base64-encoded thumbnail image data (populated by caller).
// This is not fetched automatically - use GetProtectEventThumbnail to fetch it.
type ThumbnailData struct {
	Base64   string `json:"base64,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

// ThumbnailBase64 holds the thumbnail data if fetched (not part of API response).
var _ = ThumbnailData{} // Ensure type is used

// MarshalJSON customizes JSON output to include the formatted message.
func (p *ProtectLogEntry) MarshalJSON() ([]byte, error) {
	type Alias ProtectLogEntry

	return json.Marshal(&struct {
		*Alias
		Message         string `json:"message"`
		EventType       string `json:"event_type,omitempty"`
		Category        string `json:"category_name,omitempty"`
		Severity        string `json:"severity,omitempty"`
		Title           string `json:"title,omitempty"`
		CameraName      string `json:"camera_name,omitempty"`
		UserName        string `json:"user_name,omitempty"`
		SourceName      string `json:"source_name,omitempty"`
		ThumbnailBase64 string `json:"thumbnail_base64,omitempty"`
	}{
		Alias:           (*Alias)(p),
		Message:         p.Msg(),
		EventType:       p.GetEventType(),
		Category:        p.GetCategory(),
		Severity:        p.GetSeverity(),
		Title:           p.GetTitle(),
		CameraName:      p.GetCameraName(),
		UserName:        p.GetUserName(),
		SourceName:      p.SourceName,
		ThumbnailBase64: p.ThumbnailBase64,
	})
}

// GetProtectEventThumbnail fetches the thumbnail image for a Protect event.
// Returns the raw image bytes (typically JPEG).
func (u *Unifi) GetProtectEventThumbnail(eventID string) ([]byte, error) {
	if eventID == "" {
		return nil, fmt.Errorf("event ID is required")
	}

	path := fmt.Sprintf("%s/%s/thumbnail", APIProtectEventsPath, eventID)

	u.DebugLog("Fetching Protect event thumbnail: %s", eventID)

	req, err := u.UniReq(path, "")
	if err != nil {
		return nil, fmt.Errorf("creating thumbnail request: %w", err)
	}

	// Override Accept header for binary data
	req.Header.Set("Accept", "image/jpeg, image/*")

	// Apply context timeout (same pattern as u.do() for other requests)
	var (
		cancel func()
		ctx    = context.Background()
	)

	if u.Config.Timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, u.Config.Timeout)
		defer cancel()
	}

	resp, err := u.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("fetching thumbnail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("thumbnail request failed: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading thumbnail data: %w", err)
	}

	return data, nil
}
