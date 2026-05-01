package unifi

import "fmt"

// WANStatus represents the WAN interface status from /api/s/{site}/stat/status.
type WANStatus struct {
	WANInterfaces []WANStatusInterface `json:"wan_interfaces"`

	SiteName string `json:"-"`
}

// WANStatusInterface represents a single WAN interface in the status response.
type WANStatusInterface struct {
	Name            string `fake:"{lexify:wan?}"                               json:"name"`
	State           string `fake:"{randomstring:[ACTIVE,BACKUP,DISCONNECTED]}" json:"state"`
	WANNetworkgroup string `fake:"{randomstring:[WAN,WAN2]}"                   json:"wan_networkgroup"`
}

// GetWANStatus returns the WAN interface status for a single site.
// Uses the legacy API endpoint: GET /api/s/{site}/stat/status.
// When no WAN data is returned (e.g. site has no gateway), a zero-value WANStatus with an
// empty WANInterfaces slice is returned. Callers can detect this by checking len(status.WANInterfaces) == 0.
func (u *Unifi) GetWANStatus(site *Site) (*WANStatus, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for WAN status, site %s", site.Name)

	path := fmt.Sprintf(APIWANStatusPath, site.Name)

	var response struct {
		Data []WANStatus `json:"data"`
	}

	if err := u.GetData(path, &response); err != nil {
		return nil, fmt.Errorf("fetching WAN status for site %s: %w", site.Name, err)
	}

	if len(response.Data) == 0 {
		u.DebugLog("No WAN status found for site %s", site.Name)

		return &WANStatus{SiteName: site.SiteName}, nil
	}

	if len(response.Data) > 1 {
		u.DebugLog("WAN status response for site %s contained %d entries; using first", site.Name, len(response.Data))
	}

	response.Data[0].SiteName = site.SiteName

	return &response.Data[0], nil
}
