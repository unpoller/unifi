package unifi

import "fmt"

// GetACLRules returns access control rules for a site.
func (u *Unifi) GetACLRules(site *IntegrationSite) ([]*ACLRule, error) {
	if site == nil {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Integration/v1 for ACL rules, site %s", site.Name)

	path := fmt.Sprintf(APIACLRulesPath, site.ID)

	items, err := getIntegrationList[ACLRule](u, path)
	if err != nil {
		return nil, fmt.Errorf("fetching ACL rules for site %s: %w", site.Name, err)
	}

	result := make([]*ACLRule, len(items))

	for i := range items {
		items[i].SiteName = site.Name
		result[i] = &items[i]
	}

	return result, nil
}
