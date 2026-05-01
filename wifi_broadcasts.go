package unifi

import "fmt"

// GetWifiBroadcasts returns WiFi broadcast (SSID) configurations for a site.
func (u *Unifi) GetWifiBroadcasts(site *IntegrationSite) ([]*WifiBroadcast, error) {
	if site == nil {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Integration/v1 for WiFi broadcasts, site %s", site.Name)

	path := fmt.Sprintf(APIWifiBroadcastsPath, site.ID)

	items, err := getIntegrationList[WifiBroadcast](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching wifi broadcasts for site %s: %w", site.Name, err)
	}

	result := make([]*WifiBroadcast, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
