package unifi

import (
	"encoding/json"
	"fmt"
)

// ErrInvalidMinutes is returned when AuthorizeGuest is called with a negative duration.
var ErrInvalidMinutes = fmt.Errorf("minutes must be >= 0")

// ErrEmptyMAC is returned when AuthorizeGuest is called with an empty MAC.
var ErrEmptyMAC = fmt.Errorf("mac must not be empty")

// AuthorizeGuest authorizes a guest client (identified by MAC address) on the
// given site via the stamgr command endpoint. minutes sets the authorization
// duration; pass 0 to use the controller's default. This wraps a POST to
// /api/s/{site}/cmd/stamgr with cmd=authorize-guest.
func (u *Unifi) AuthorizeGuest(site *Site, mac string, minutes int) error {
	if site == nil || site.Name == "" {
		return ErrNoSiteProvided
	}

	if mac == "" {
		return ErrEmptyMAC
	}

	if minutes < 0 {
		return ErrInvalidMinutes
	}

	payload := struct {
		Cmd     string `json:"cmd"`
		MAC     string `json:"mac"`
		Minutes int    `json:"minutes,omitempty"`
	}{
		Cmd:     "authorize-guest",
		MAC:     mac,
		Minutes: minutes,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling authorize-guest payload: %w", err)
	}

	path := fmt.Sprintf(APIStaMgrPath, site.Name)

	if _, err := u.PostJSON(path, string(body)); err != nil {
		return fmt.Errorf("authorize-guest %s: %w", mac, err)
	}

	return nil
}
