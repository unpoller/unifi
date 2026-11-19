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

// ConsoleSystem models the verified firmware and application update fields of the
// local UniFi OS system response. The endpoint is console-wide, not scoped to a
// Network site. Other fields are ignored until their response shapes and semantics
// are validated. Firmware and application entries may be absent in a partial response.
type ConsoleSystem struct {
	SourceName string `json:"-"`
	Apps       struct {
		Controllers []*ConsoleApplication `json:"controllers"`
	} `json:"apps"`
	Devices struct {
		UniFiOS []*ConsoleFirmware `json:"unifiOS"`
	} `json:"devices"`
}

// ConsoleApplication describes installed and offered application versions from
// apps.controllers, independently of console OS firmware and Network sysinfo flags.
// An empty UpdateAvailable does not establish that the application is up to date.
type ConsoleApplication struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	UpdateAvailable string `json:"updateAvailable"`
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
