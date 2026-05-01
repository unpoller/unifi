package unifi

import "fmt"

// PortForward represents a port forwarding rule from /api/s/{site}/rest/portforward.
type PortForward struct {
	ID      string   `fake:"{uuid}"                           json:"_id"`
	Enabled FlexBool `json:"enabled"`
	FwdIP   string   `fake:"{ipv4address}"                    json:"fwd"`
	FwdPort string   `fake:"{number:1,65535}"                 json:"fwd_port"`
	Name    string   `fake:"{buzzword}"                       json:"name"`
	PfwdPf  string   `json:"pfwd_interface"`
	Proto   string   `fake:"{randomstring:[tcp,udp,tcp_udp]}" json:"proto"`
	SiteID  string   `fake:"{uuid}"                           json:"site_id"`
	Src     string   `json:"src"`
	DstPort string   `fake:"{number:1,65535}"                 json:"dst_port"`
	Log     FlexBool `json:"log"`

	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}

// GetPortForwards returns port forwarding rules for a single site.
// Uses the legacy API endpoint: GET /api/s/{site}/rest/portforward.
func (u *Unifi) GetPortForwards(site *Site) ([]*PortForward, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for port forwards, site %s", site.SiteName)

	path := fmt.Sprintf(APIPortForwardPath, site.Name)

	var response struct {
		Data []PortForward `json:"data"`
	}

	if err := u.GetData(path, &response); err != nil {
		return nil, fmt.Errorf("fetching port forwards for site %s: %w", site.SiteName, err)
	}

	result := make([]*PortForward, len(response.Data))

	for i := range response.Data {
		response.Data[i].SiteName = site.SiteName
		response.Data[i].SourceName = u.URL
		result[i] = &response.Data[i]
	}

	return result, nil
}
