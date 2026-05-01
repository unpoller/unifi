package unifi

import "fmt"

// GetHotspotVouchers returns guest portal vouchers for a site.
func (u *Unifi) GetHotspotVouchers(site *IntegrationSite) ([]*HotspotVoucher, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration/v1 for hotspot vouchers, site %s", site.Name)

	path := fmt.Sprintf(APIHotspotVouchersPath, site.ID)

	items, err := getIntegrationList[HotspotVoucher](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching hotspot vouchers for site %s: %w", site.Name, err)
	}

	result := make([]*HotspotVoucher, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
