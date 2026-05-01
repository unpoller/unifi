package unifi

import "fmt"

// GetLAGs returns link aggregation group configurations for a site.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetLAGs(site *IntegrationSite) ([]*LAG, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration API for LAG data, site %s", site.Name)

	path := fmt.Sprintf(APILAGsPath, site.ID)

	items, err := getIntegrationList[LAG](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching LAGs for site %s: %w", site.Name, err)
	}

	result := make([]*LAG, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}

// GetMCLAGDomains returns multi-chassis LAG domain configurations for a site.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetMCLAGDomains(site *IntegrationSite) ([]*MCLAGDomain, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration API for MCLAG domain data, site %s", site.Name)

	path := fmt.Sprintf(APIMCLAGDomainsPath, site.ID)

	items, err := getIntegrationList[MCLAGDomain](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching MCLAG domains for site %s: %w", site.Name, err)
	}

	result := make([]*MCLAGDomain, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}

// GetSwitchStacks returns switch stack configurations for a site.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetSwitchStacks(site *IntegrationSite) ([]*SwitchStack, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration API for switch stack data, site %s", site.Name)

	path := fmt.Sprintf(APISwitchStacksPath, site.ID)

	items, err := getIntegrationList[SwitchStack](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching switch stacks for site %s: %w", site.Name, err)
	}

	result := make([]*SwitchStack, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
