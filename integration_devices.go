package unifi

// GetIntegrationDeviceStats returns statistics for a single device from the Integration/v1 API.
func (u *Unifi) GetIntegrationDeviceStats(_ *IntegrationSite, _ string) (*IntegrationDeviceStats, error) {
	return nil, nil
}

// GetAllIntegrationDeviceStats returns statistics for all devices in a site.
func (u *Unifi) GetAllIntegrationDeviceStats(_ *IntegrationSite) ([]*IntegrationDeviceStats, error) {
	return nil, nil
}
