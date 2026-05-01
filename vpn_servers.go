package unifi

// GetVPNServers returns VPN server configurations for a site.
func (u *Unifi) GetVPNServers(_ *IntegrationSite) ([]*VPNServer, error) {
	return nil, nil
}

// GetSiteToSiteTunnels returns site-to-site VPN tunnel configurations for a site.
func (u *Unifi) GetSiteToSiteTunnels(_ *IntegrationSite) ([]*SiteToSiteTunnel, error) {
	return nil, nil
}
