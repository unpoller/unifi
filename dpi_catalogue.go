package unifi

import "fmt"

// GetDPIApplications returns the DPI application reference catalogue (global, no site).
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetDPIApplications() ([]*DPIApplication, error) {
	items, err := getIntegrationList[DPIApplication](u, APIDPIApplicationsPath)
	if err != nil {
		return nil, fmt.Errorf("fetching DPI applications: %w", err)
	}

	result := make([]*DPIApplication, len(items))

	for i := range items {
		result[i] = &items[i]
	}

	return result, nil
}

// GetDPICategories returns the DPI category reference catalogue (global, no site).
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetDPICategories() ([]*DPICategory, error) {
	items, err := getIntegrationList[DPICategory](u, APIDPICategoriesPath)
	if err != nil {
		return nil, fmt.Errorf("fetching DPI categories: %w", err)
	}

	result := make([]*DPICategory, len(items))

	for i := range items {
		result[i] = &items[i]
	}

	return result, nil
}
