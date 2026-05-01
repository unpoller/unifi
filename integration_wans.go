package unifi

import "fmt"

// GetIntegrationWANs returns WAN interface identifiers for a site from the Integration/v1 API.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetIntegrationWANs(site *IntegrationSite) ([]*IntegrationWAN, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration/v1 for WAN interfaces, site %s", site.Name)

	path := fmt.Sprintf(APIIntegrationWANsPath, site.ID)

	items, err := getIntegrationList[IntegrationWAN](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching integration WANs for site %s: %w", site.Name, err)
	}

	result := make([]*IntegrationWAN, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
