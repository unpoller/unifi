package unifi

import "fmt"

// GetIntegrationNetworks returns networks for a site from the Integration/v1 API.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetIntegrationNetworks(site *IntegrationSite) ([]*IntegrationNetwork, error) {
	if site == nil {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Integration/v1 for networks, site %s", site.Name)

	path := fmt.Sprintf(APIIntegrationNetworksPath, site.ID)

	items, err := getIntegrationList[IntegrationNetwork](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching integration networks for site %s: %w", site.Name, err)
	}

	result := make([]*IntegrationNetwork, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
