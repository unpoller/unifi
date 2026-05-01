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
// Returns nil, nil for an empty result set (totalCount=0), following Go convention.
func getIntegrationList[T any](u *Unifi, path string) ([]T, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if u.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	u.DebugLog("Polling Integration/v1 list %s", path)

	const pageSize = 200

	var (
		all             []T
		totalCount      int
		driftCount      int // incremented for every page where totalCount changes from the first-page value
		driftOffset     int // first offset where drift was recorded; always >0 (drift only occurs in else-if when offset!=0); 0 means no drift
		driftTotalCount int // the totalCount value the server reported at driftOffset
	)

	offset := 0

	// The loop is bounded: offset advances by page.Count each iteration and exits when
	// offset >= totalCount. totalCount is captured from the first page; if the server grows
	// totalCount after a full first page is already returned, the loop exits without a second
	// page to compare against and drift goes undetected — inherent in snapshot-based pagination.
	// The page.Count == 0 guard handles servers that return an empty page prematurely.
	for {
		pagedPath := fmt.Sprintf("%s?offset=%d&limit=%d", path, offset, pageSize)

		body, err := u.GetJSON(pagedPath)
		if err != nil {
			u.ErrorLog("integration page %s: request failed: %v", pagedPath, err)

			return nil, fmt.Errorf("fetching integration page %s: %w", pagedPath, err)
		}

		u.DebugLog("Fetched integration page %s (%d bytes)", pagedPath, len(body))

		var page integrationPage[T]

		if err := json.Unmarshal(body, &page); err != nil {
			// DebugLog rather than ErrorLog: the returned error is the actionable signal for callers;
			// the body is supplementary context (may be a large HTML proxy-error page) only needed
			// when actively debugging. %.512s is a byte-boundary truncation and may split a multi-byte
			// UTF-8 sequence, but JSON responses are ASCII-safe enough for diagnostic use.
			u.DebugLog("integration list %s: failed response body (first 512 bytes): %.512s", pagedPath, body)

			return nil, fmt.Errorf("parsing integration response from %s: %w", pagedPath, err)
		}

		// A per-page count mismatch is a protocol violation distinct from the multi-page
		// completeness failure that ErrIncompleteResults signals: the server contradicted
		// itself within a single response rather than across pages. Not wrapping
		// ErrIncompleteResults is intentional — callers distinguish the two modes by
		// checking errors.Is(err, ErrIncompleteResults): false means per-page violation,
		// true means multi-page count mismatch.
		if page.Count != len(page.Data) {
			u.ErrorLog("integration page %s: server reported count=%d but returned %d items", pagedPath, page.Count, len(page.Data))

			return nil, fmt.Errorf("integration page %s: server reported count=%d but returned %d items", pagedPath, page.Count, len(page.Data))
		}

		// Check for invalid totalCount before drift detection: a negative value on a later
		// page must return an error rather than record a spurious driftOffset.
		if page.TotalCount < 0 {
			u.ErrorLog("integration page %s: server returned negative totalCount=%d", pagedPath, page.TotalCount)

			return nil, fmt.Errorf("integration page %s: server returned negative totalCount=%d", pagedPath, page.TotalCount)
		}

		if offset == 0 {
			totalCount = page.TotalCount
		} else if page.TotalCount != totalCount {
			driftCount++

			if driftOffset == 0 {
				driftOffset = offset
				driftTotalCount = page.TotalCount // logged at end with final counts for full context
			}
		}

		all = append(all, page.Data...)
		offset += page.Count

		if page.Count == 0 {
			// "integration list" (path) rather than "integration page" (pagedPath) is intentional:
			// this is a list-level event; the offset field already identifies the triggering page.
			// ErrIncompleteResults covers both early-exit and final-count-mismatch paths; callers
			// that need to distinguish them can inspect the error message from the preceding ErrorLog.
			if offset < totalCount {
				u.ErrorLog("integration list %s: server returned count=0 at offset %d but only %d/%d items fetched; stopping early", path, offset, len(all), totalCount)

				return nil, ErrIncompleteResults
			}

			u.DebugLog("integration list %s: server returned count=0 with totalCount=%d; done", path, totalCount)

			break
		}

		if offset >= totalCount {
			break
		}
	}

	if len(all) != totalCount {
		if driftCount > 0 {
			u.ErrorLog("integration list %s: totalCount drifted %d time(s), first at offset %d from %d to %d; fetched %d items", path, driftCount, driftOffset, totalCount, driftTotalCount, len(all))
		} else {
			u.ErrorLog("integration list %s: fetched %d items but server reported totalCount=%d", path, len(all), totalCount)
		}

		return nil, ErrIncompleteResults
	}

	if driftCount > 0 {
		// DebugLog rather than ErrorLog: all items were fetched successfully; drift that resolved
		// is informational, not a failure. ErrorLog here would produce false production alerts.
		u.DebugLog("integration list %s: totalCount drifted %d time(s) but resolved; first drift at offset %d from %d to %d", path, driftCount, driftOffset, totalCount, driftTotalCount)
	}

	u.DebugLog("Fetched %d/%d items from %s", len(all), totalCount, path)

	return all, nil
}
