package unifi

import "fmt"

// GetTrafficMatchingLists returns traffic matching lists for a site.
func (u *Unifi) GetTrafficMatchingLists(site *IntegrationSite) ([]*TrafficMatchingList, error) {
	if site == nil {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Integration/v1 for traffic matching lists, site %s", site.Name)

	path := fmt.Sprintf(APITrafficMatchingListsPath, site.ID)

	items, err := getIntegrationList[TrafficMatchingList](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching traffic matching lists for site %s: %w", site.Name, err)
	}

	result := make([]*TrafficMatchingList, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
