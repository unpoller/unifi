package unifi

import "fmt"

// GetRADIUSProfiles returns RADIUS authentication profiles for a site.
func (u *Unifi) GetRADIUSProfiles(site *IntegrationSite) ([]*RADIUSProfile, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil {
		return nil, ErrNoSiteProvided
	}

	if site.ID == "" {
		return nil, fmt.Errorf("site %q has an empty ID; cannot construct Integration/v1 API path", site.Name)
	}

	u.DebugLog("Polling Integration/v1 for RADIUS profiles, site %s", site.Name)

	path := fmt.Sprintf(APIRADIUSProfilesPath, site.ID)

	items, err := getIntegrationList[RADIUSProfile](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching RADIUS profiles for site %s: %w", site.Name, err)
	}

	result := make([]*RADIUSProfile, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
