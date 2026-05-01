package unifi

import "fmt"

// GetDNSPolicies returns DNS policies for a site.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetDNSPolicies(site *IntegrationSite) ([]*DNSPolicy, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration API for DNS policy data, site %s", site.Name)

	path := fmt.Sprintf(APIDNSPoliciesPath, site.ID)

	items, err := getIntegrationList[DNSPolicy](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching DNS policies for site %s: %w", site.Name, err)
	}

	result := make([]*DNSPolicy, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
