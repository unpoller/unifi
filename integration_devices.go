package unifi

import (
	"encoding/json"
	"fmt"
)

// GetIntegrationDeviceStats returns statistics for a single device from the Integration/v1 API.
func (u *Unifi) GetIntegrationDeviceStats(site *IntegrationSite, deviceID string) (*IntegrationDeviceStats, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	if deviceID == "" {
		return nil, fmt.Errorf("deviceID must not be empty")
	}

	u.DebugLog("Polling Integration/v1 for device stats, site %s device %s", site.Name, deviceID)

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
// If fetching stats for a device fails, partial results collected so far are returned alongside the error.
func (u *Unifi) GetAllIntegrationDeviceStats(site *IntegrationSite) ([]*IntegrationDeviceStats, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration/v1 for all device stats, site %s", site.Name)

	type integrationDeviceID struct {
		ID string `json:"id"`
	}

	devPath := fmt.Sprintf(APIIntegrationDevicesPath, site.ID)

	devices, err := getIntegrationList[integrationDeviceID](u, devPath)
	if err != nil {
		return nil, fmt.Errorf("fetching device list for site %s: %w", site.Name, err)
	}

	result := make([]*IntegrationDeviceStats, 0, len(devices))

	skipped := 0

	for _, dev := range devices {
		if dev.ID == "" {
			skipped++

			continue
		}

		stats, err := u.GetIntegrationDeviceStats(site, dev.ID)
		if err != nil {
			// DebugLog provides supplementary context; the returned error is the primary signal for the caller.
			if skipped > 0 {
				u.DebugLog("Skipped %d/%d devices with empty IDs in site %s before error", skipped, len(devices), site.Name)
			}

			return result, fmt.Errorf("fetching stats for device %s in site %s: %w", dev.ID, site.Name, err)
		}

		result = append(result, stats)
	}

	// DebugLog rather than ErrorLog: an empty device ID is a server data quality issue (e.g.,
	// a device mid-adoption), not a code error. ErrorLog would produce false production alerts.
	if skipped > 0 {
		u.DebugLog("Skipped %d/%d devices with empty IDs in site %s", skipped, len(devices), site.Name)
	}

	return result, nil
}
