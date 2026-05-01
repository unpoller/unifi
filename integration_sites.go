package unifi

import (
	"encoding/json"
	"fmt"
)

// GetIntegrationSites returns all sites from the Integration/v1 API.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
// Use InternalReference to cross-reference with the legacy Site.Name field.
func (u *Unifi) GetIntegrationSites() ([]*IntegrationSite, error) {
	sites, err := getIntegrationList[IntegrationSite](u, APIIntegrationSitesPath)
	if err != nil {
		return nil, fmt.Errorf("fetching integration sites: %w", err)
	}

	result := make([]*IntegrationSite, len(sites))

	for i := range sites {
		result[i] = &sites[i]
	}

	return result, nil
}

// GetIntegrationInfo returns application version info from the Integration/v1 API.
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetIntegrationInfo() (*IntegrationInfo, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	body, err := u.GetJSON(APIIntegrationInfoPath)
	if err != nil {
		return nil, fmt.Errorf("fetching integration info: %w", err)
	}

	var info IntegrationInfo

	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("parsing integration info: %w", err)
	}

	return &info, nil
}
