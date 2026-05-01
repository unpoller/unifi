package unifi

import "fmt"

// GetCountries returns the list of countries for geo-based firewall filters (global, no site).
// Requires Config.APIKey; returns ErrAPIKeyRequired when no key is configured.
func (u *Unifi) GetCountries() ([]*Country, error) {
	items, err := getIntegrationList[Country](u, APICountriesPath)
	if err != nil {
		return nil, fmt.Errorf("fetching countries: %w", err)
	}

	result := make([]*Country, len(items))

	for i := range items {
		result[i] = &items[i]
	}

	return result, nil
}
