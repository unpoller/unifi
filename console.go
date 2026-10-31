package unifi

import (
	"errors"
	"fmt"
)

// APIConsoleSystemPath is served by UniFi OS itself, outside the Network proxy.
const APIConsoleSystemPath = "/api/system"

var (
	// ErrConsoleSystemUnsupported indicates a controller without UniFi OS.
	ErrConsoleSystemUnsupported = errors.New("console system data requires UniFi OS")
	// ErrNoConsoleSystemData indicates a response with no console firmware entries.
	ErrNoConsoleSystemData = errors.New("console system response contains no firmware data")
)

// ConsoleSystem contains console-wide information from the local UniFi OS API.
// Application updates are intentionally not included: their response shape needs
// separate validation. This endpoint is not scoped to a Network site.
type ConsoleSystem struct {
	SourceName string `json:"-"`
	Devices    struct {
		UniFiOS []*ConsoleFirmware `json:"unifiOS"`
	} `json:"devices"`
}

// ConsoleFirmware describes the installed and offered UniFi OS firmware.
type ConsoleFirmware struct {
	Name            string `json:"name"`
	Firmware        string `json:"firmware"`
	FirmwareStatus  string `json:"firmwareStatus"`
	UpdateAvailable string `json:"updateAvailable"`
}

// UpdateStatus reports whether an update is available and whether that state is
// known. A missing/null offered version alone does not establish "no update".
// Unrecognized firmware states without an offered version remain unknown.
func (f *ConsoleFirmware) UpdateStatus() (available, known bool) {
	if f == nil {
		return false, false
	}

	if f.UpdateAvailable != "" || f.FirmwareStatus == "updateAvailable" {
		return true, true
	}

	if f.FirmwareStatus == "upToDate" {
		return false, true
	}

	return false, false
}

// GetConsoleSystem reads the local console's firmware state using the existing
// authenticated session. Permission and unsupported-endpoint errors are returned
// to the caller; they must not be interpreted as an up-to-date console.
// Like the Protect getters, this is not part of the Network UnifiClient interface.
func (u *Unifi) GetConsoleSystem() (*ConsoleSystem, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if !u.new {
		return nil, ErrConsoleSystemUnsupported
	}

	var response ConsoleSystem

	if err := u.GetData(APIConsoleSystemPath, &response); err != nil {
		return nil, fmt.Errorf("getting console system: %w", err)
	}

	for _, firmware := range response.Devices.UniFiOS {
		if firmware != nil {
			response.SourceName = u.URL

			return &response, nil
		}
	}

	return nil, ErrNoConsoleSystemData
}
