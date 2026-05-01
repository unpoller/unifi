package unifi

import "fmt"

// GetPendingDevices returns devices waiting to be adopted (global, no site).
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetPendingDevices() ([]*PendingDevice, error) {
	items, err := getIntegrationList[PendingDevice](u, APIPendingDevicesPath)
	if err != nil {
		return nil, fmt.Errorf("fetching pending devices: %w", err)
	}

	result := make([]*PendingDevice, len(items))

	for i := range items {
		result[i] = &items[i]
	}

	return result, nil
}
