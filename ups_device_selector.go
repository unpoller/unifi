package unifi

import "fmt"

// UPSDeviceSelector represents a UPS device entry from /api/s/{site}/stat/ups-devices.
// This is a lightweight selector format distinct from the full device data in /stat/device.
type UPSDeviceSelector struct {
	ID     string `fake:"{uuid}"        json:"_id"`
	Image  string `fake:"{imageurl}"    json:"image"`
	Label  string `fake:"{productname}" json:"label"`
	MAC    string `fake:"{macaddress}"  json:"mac"`
	SiteID string `fake:"{uuid}"        json:"site_id"`

	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}

// GetUPSDeviceList returns the list of UPS device selectors for a single site.
// Uses the legacy API endpoint: GET /api/s/{site}/stat/ups-devices.
func (u *Unifi) GetUPSDeviceList(site *Site) ([]*UPSDeviceSelector, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for UPS device list, site %s", site.Name)

	path := fmt.Sprintf(APIUPSDevicesPath, site.Name)

	var response struct {
		Data []UPSDeviceSelector `json:"data"`
	}

	if err := u.GetData(path, &response); err != nil {
		return nil, fmt.Errorf("fetching UPS device list for site %s: %w", site.Name, err)
	}

	result := make([]*UPSDeviceSelector, len(response.Data))

	for i := range response.Data {
		response.Data[i].SiteName = site.SiteName
		response.Data[i].SourceName = u.URL
		result[i] = &response.Data[i]
	}

	return result, nil
}
