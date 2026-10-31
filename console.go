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
	// ErrNoConsoleSystemData indicates a null console system response.
	ErrNoConsoleSystemData = errors.New("console system response is null")
)

// ConsoleSystem models the verified firmware subset of the local UniFi OS system
// response. The endpoint is console-wide, not scoped to a Network site. Other
// fields are ignored until their response shapes and semantics are validated.
// Firmware entries may be absent in a valid partial response.
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

// GetConsoleSystem reads the modeled fields from the local console system response
// using the existing authenticated session. Valid partial objects are accepted;
// absent firmware data does not establish an up-to-date console. Permission and
// unsupported-endpoint errors are returned to the caller.
// Like the Protect getters, this is not part of the Network UnifiClient interface.
func (u *Unifi) GetConsoleSystem() (*ConsoleSystem, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if !u.new {
		return nil, ErrConsoleSystemUnsupported
	}

	var response *ConsoleSystem

	if err := u.GetData(APIConsoleSystemPath, &response); err != nil {
		return nil, fmt.Errorf("getting console system: %w", err)
	}

	if response == nil {
		return nil, ErrNoConsoleSystemData
	}

	response.SourceName = u.URL

	return response, nil
}
