package unifi

import (
	"encoding/json"
	"fmt"
)

// GetIntegrationDeviceStats returns statistics for a single device from the Integration/v1 API.
func (u *Unifi) GetIntegrationDeviceStats(site *IntegrationSite, deviceID string) (*IntegrationDeviceStats, error) {
	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	path := fmt.Sprintf(APIIntegrationDeviceStatsPath, site.ID, deviceID)

	body, err := u.GetJSON(path)
	if err != nil {
		return nil, fmt.Errorf("fetching device stats for site %s device %s: %w", site.Name, deviceID, err)
	}

	var stats IntegrationDeviceStats

	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("parsing device stats for site %s device %s: %w", site.Name, deviceID, err)
	}

	return &stats, nil
}

// GetAllIntegrationDeviceStats returns statistics for all devices in a site.
func (u *Unifi) GetAllIntegrationDeviceStats(site *IntegrationSite) ([]*IntegrationDeviceStats, error) {
	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	type integrationDeviceID struct {
		ID string `json:"id"`
	}

	devPath := fmt.Sprintf(APIIntegrationDevicesPath, site.ID)

	devices, err := getIntegrationList[integrationDeviceID](u, devPath)
	if err != nil {
		return nil, fmt.Errorf("fetching device list for site %s: %w", site.Name, err)
	}

	result := make([]*IntegrationDeviceStats, 0, len(devices))

	for _, dev := range devices {
		stats, err := u.GetIntegrationDeviceStats(site, dev.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, stats)
	}

	return result, nil
}
