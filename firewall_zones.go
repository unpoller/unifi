package unifi

import "fmt"

// GetFirewallZones returns firewall zones for a site. Zone IDs appear in firewall policies.
func (u *Unifi) GetFirewallZones(site *IntegrationSite) ([]*FirewallZone, error) {
	if site == nil {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Integration/v1 for firewall zones, site %s", site.Name)

	path := fmt.Sprintf(APIFirewallZonesPath, site.ID)

	items, err := getIntegrationList[FirewallZone](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching firewall zones for site %s: %w", site.Name, err)
	}

	result := make([]*FirewallZone, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
