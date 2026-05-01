package unifi

// GetLAGs returns link aggregation group configurations for a site.
func (u *Unifi) GetLAGs(_ *IntegrationSite) ([]*LAG, error) {
	return nil, nil
}

// GetMCLAGDomains returns multi-chassis LAG domain configurations for a site.
func (u *Unifi) GetMCLAGDomains(_ *IntegrationSite) ([]*MCLAGDomain, error) {
	return nil, nil
}

// GetSwitchStacks returns switch stack configurations for a site.
func (u *Unifi) GetSwitchStacks(_ *IntegrationSite) ([]*SwitchStack, error) {
	return nil, nil
}
