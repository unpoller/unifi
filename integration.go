package unifi

import (
	"encoding/json"
	"fmt"
)

// integrationPage is the JSON envelope for Integration/v1 paginated list responses.
type integrationPage[T any] struct {
	Count      int `json:"count"`
	Data       []T `json:"data"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	TotalCount int `json:"totalCount"`
}

// getIntegrationList fetches all pages from an Integration/v1 list endpoint.
// path must not include offset/limit query parameters.
func getIntegrationList[T any](u *Unifi, path string) ([]T, error) {
	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	const pageSize = 200

	var all []T

	offset := 0

	for {
		pagedPath := fmt.Sprintf("%s?offset=%d&limit=%d", path, offset, pageSize)

		body, err := u.GetJSON(pagedPath)
		if err != nil {
			return nil, err
		}

		var page integrationPage[T]

		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("parsing integration response: %w", err)
		}

		all = append(all, page.Data...)
		offset += page.Count

		if offset >= page.TotalCount || page.Count == 0 {
			break
		}
	}

	return all, nil
}
