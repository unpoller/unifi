package unifi

import "fmt"

// GetVPNServers returns VPN server configurations for a site from the Integration/v1 API.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetVPNServers(site *IntegrationSite) ([]*VPNServer, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration/v1 for VPN servers, site %s", site.Name)

	path := fmt.Sprintf(APIVPNServersPath, site.ID)

	items, err := getIntegrationList[VPNServer](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching VPN servers for site %s: %w", site.Name, err)
	}

	result := make([]*VPNServer, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}

// GetSiteToSiteTunnels returns site-to-site VPN tunnel configurations for a site from the Integration/v1 API.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetSiteToSiteTunnels(site *IntegrationSite) ([]*SiteToSiteTunnel, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration/v1 for site-to-site tunnels, site %s", site.Name)

	path := fmt.Sprintf(APISiteToSiteTunnelsPath, site.ID)

	items, err := getIntegrationList[SiteToSiteTunnel](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching site-to-site tunnels for site %s: %w", site.Name, err)
	}

	result := make([]*SiteToSiteTunnel, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
