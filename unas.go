package unifi

import (
	"crypto/sha256"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// UNAS Pro support is a port of https://github.com/alexgreenbank/unaspoller (MIT), which
// worked out the UniFi Drive JSON API by observing the console's own web UI. Thanks to
// alexgreenbank for the reverse engineering and for offering the work upstream; see
// unpoller/unpoller#785.

// unasEndpointCount is the number of endpoints GetUNASDevice polls. It is used to tell a
// partial failure (report what we have) from a total one (report nothing).
const unasEndpointCount = 4

// UNASDeviceInfo holds console identity, CPU, and memory data.
// API Path: /proxy/drive/api/v2/systems/device-info
type UNASDeviceInfo struct {
	Name              string                 `json:"name"`
	Model             string                 `json:"model"`
	Version           string                 `json:"version"`
	FirmwareVersion   string                 `json:"firmwareVersion"`
	Status            string                 `json:"status"`
	StartupTime       string                 `json:"startupTime"`
	SfpAggregation    FlexBool               `json:"sfpAggregation"`
	Memory            UNASMemory             `json:"memory"`
	CPU               UNASCPU                `json:"cpu"`
	NetworkInterfaces []UNASNetworkInterface `json:"networkInterfaces"`
	// Usbs is nullable and its populated shape is unknown; kept raw for diagnosis.
	Usbs json.RawMessage `fake:"skip" json:"usbs"`
}

// UNASMemory holds console memory counters, in bytes.
type UNASMemory struct {
	Free      FlexInt `json:"free"`
	Total     FlexInt `json:"total"`
	Available FlexInt `json:"available"`
}

// UNASCPU holds console CPU load and temperature.
type UNASCPU struct {
	CurrentLoad FlexInt `json:"currentLoad"`
	Temperature FlexInt `json:"temperature"`
}

// UNASNetworkInterface describes one console network interface.
type UNASNetworkInterface struct {
	Interface     string   `json:"interface"`
	InterfaceName string   `json:"interfaceName"`
	Connected     FlexBool `json:"connected"`
	MaxSpeed      string   `json:"maxSpeed"`
	LinkSpeed     string   `json:"linkSpeed"`
	Address       string   `json:"address,omitempty"`
	MAC           string   `json:"mac,omitempty"`
}

// UNASStorage holds storage pool and physical disk data.
// API Path: /proxy/drive/api/v2/storage
type UNASStorage struct {
	Pools []UNASPool `json:"pools"`
	Disks []UNASDisk `json:"disks"`
	// CacheSlots and Expansions have no known populated shape yet; kept raw for diagnosis.
	CacheSlots json.RawMessage `fake:"skip" json:"cacheSlots"`
	Expansions json.RawMessage `fake:"skip" json:"expansions"`
}

// UNASPool describes one storage pool.
type UNASPool struct {
	Number             FlexInt         `json:"number"`
	ID                 string          `json:"id"`
	PreferLevel        string          `json:"preferLevel"`
	Type               string          `json:"type"`
	Status             string          `json:"status"`
	Capacity           FlexInt         `json:"capacity"`
	Usage              FlexInt         `json:"usage"`
	ActiveRaidGroupID  string          `json:"activeRaidGroupId"`
	InitializingStatus string          `json:"initializingStatus"`
	RaidGroups         []UNASRaidGroup `json:"raidGroups"`
}

// UNASRaidGroup describes one RAID group within a storage pool.
type UNASRaidGroup struct {
	Number              FlexInt  `json:"number"`
	ID                  string   `json:"id"`
	RemnantReason       string   `json:"remnantReason"`
	IsSSDCache          FlexBool `json:"isSSDCache"`
	CurrentLevel        string   `json:"currentLevel"`
	ConfigLevel         string   `json:"configLevel"`
	CurrentProtection   FlexInt  `json:"currentProtection"`
	ExpectedProtection  FlexInt  `json:"expectedProtection"`
	RecommendedDiskSize FlexInt  `json:"recommendedDiskSize"`
	Progress            FlexInt  `json:"progress"`
	Estimate            FlexInt  `json:"estimate"`
}

// UNASDisk describes one physical disk, including SMART health data.
type UNASDisk struct {
	SlotID                   string   `json:"slotId"`
	Location                 string   `json:"location"`
	PoolID                   string   `json:"poolId"`
	RaidGroupID              string   `json:"raidGroupId"`
	MetadataGroupID          string   `json:"metadataGroupId"`
	IsGlobalHotSpare         FlexBool `json:"isGlobalHotSpare"`
	IsLocalHotSpare          FlexBool `json:"isLocalHotSpare"`
	Type                     string   `json:"type"`
	State                    string   `json:"state"`
	RPM                      FlexInt  `json:"rpm"`
	Model                    string   `json:"model"`
	Size                     FlexInt  `json:"size"`
	Sata                     string   `json:"sata"`
	Ata                      string   `json:"ata"`
	NvmeVersion              string   `json:"nvmeVersion"`
	Firmware                 string   `json:"firmware"`
	SectorFormat             string   `json:"sectorFormat"`
	Serial                   string   `json:"serial"`
	Temperature              FlexInt  `json:"temperature"`
	PowerOnHours             FlexInt  `json:"powerOnHours"`
	BadSectorCount           FlexInt  `json:"badSectorCount"`
	UncorrectableSectorCount FlexInt  `json:"uncorrectableSectorCount"`
	ReadErrorRate            FlexInt  `json:"readErrorRate"`
	SmartReadErrorCount      FlexInt  `json:"smartReadErrorCount"`
	ReadKBPS                 FlexInt  `json:"readKBPS"`
	WriteKBPS                FlexInt  `json:"writeKBPS"`
	SmartTestSupported       FlexBool `json:"smartTestSupported"`
	HealthScore              FlexInt  `json:"healthScore"`
	// RiskReasons and IncompatibleReasons have no known populated shape yet.
	RiskReasons         json.RawMessage `fake:"skip" json:"riskReasons"`
	IncompatibleReasons json.RawMessage `fake:"skip" json:"incompatibleReasons"`
}

// UNASDrive describes one share ("drive") on the console.
type UNASDrive struct {
	ID                       string               `json:"id"`
	Type                     string               `json:"type"`
	Name                     string               `json:"name"`
	Status                   string               `json:"status"`
	StoragePoolID            string               `json:"storagePoolId"`
	DataSync                 string               `json:"dataSync"`
	RecordSize               string               `json:"recordSize"`
	CompressionLevel         string               `json:"compressionLevel"`
	Deduplication            string               `json:"deduplication"`
	DeduplicationEverEnabled FlexBool             `json:"deduplicationEverEnabled"`
	Quota                    FlexInt              `json:"quota"`
	Usage                    FlexInt              `json:"usage"`
	Role                     string               `json:"role"`
	MemberCount              FlexInt              `json:"memberCount"`
	Protections              UNASDriveProtections `json:"protections"`
}

// UNASDriveProtections holds the protection settings for a share.
type UNASDriveProtections struct {
	EncryptionStatus    string   `json:"encryptionStatus"`
	RemoteBackupEnabled FlexBool `json:"remoteBackupEnabled"`
	SnapshotEnabled     FlexBool `json:"snapshotEnabled"`
}

// UNASNetworkIO holds instantaneous console network throughput.
// API Path: /proxy/drive/api/v2/systems/network-io
type UNASNetworkIO struct {
	ReceiveKBPS  FlexInt `json:"receiveKBPS"`
	TransmitKBPS FlexInt `json:"transmitKBPS"`
	Timestamp    string  `json:"timestamp"`
}

// UNASDevice is the aggregate of every polled UNAS endpoint for one console.
// Any member may be nil if that endpoint failed; check before dereferencing.
type UNASDevice struct {
	DeviceInfo *UNASDeviceInfo `json:"device_info"`
	Storage    *UNASStorage    `json:"storage"`
	NetworkIO  *UNASNetworkIO  `json:"network_io"`
	Drives     []*UNASDrive    `json:"drives"`

	SourceName string `json:"source_name"`
}

// Name returns the console name, or an empty string if device-info is missing.
func (d *UNASDevice) Name() string {
	if d == nil || d.DeviceInfo == nil {
		return ""
	}

	return d.DeviceInfo.Name
}

// Model returns the console model, or an empty string if device-info is missing.
func (d *UNASDevice) Model() string {
	if d == nil || d.DeviceInfo == nil {
		return ""
	}

	return d.DeviceInfo.Model
}

// UNASClient talks to a UNAS Pro console's UniFi Drive API.
//
// It embeds *Unifi to reuse the login, CSRF, cookie jar, gzip, and rate-limit handling,
// but it is constructed differently: see NewUNASClient. Drive getters are deliberately
// not part of the UnifiClient interface, since a UNAS console serves none of the rest of
// it and mocks.MockUnifi would have to grow methods it can never meaningfully implement.
type UNASClient struct {
	*Unifi
}

// NewUNASClient creates an authenticated client for a UNAS Pro console. Start here for UNAS.
//
// This is deliberately not NewUnifi: that constructor finishes with GetServerData(), which
// GETs /status — rewritten to /proxy/network/status — and a storage-only UNAS console runs
// no Network application to answer it. So this logs in and stops.
//
// API-key auth is rejected for the same reason: Login() uses /status as its login path when
// Config.APIKey is set.
func NewUNASClient(config *Config) (*UNASClient, error) {
	if config == nil {
		return nil, ErrNilConfig
	}

	if config.APIKey != "" {
		return nil, ErrAPIKeyUnsupported
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, fmt.Errorf("creating cookiejar: %w", err)
	}

	u := newUnifi(config, jar)

	for i, cert := range config.SSLCert {
		p, _ := pem.Decode(cert)
		if p == nil {
			continue
		}

		u.fingerprints[i] = fmt.Sprintf("%x", sha256.Sum256(p.Bytes))
	}

	// A UNAS Pro is always a UniFi OS console, so the new-style API is a given. Setting this
	// directly rather than calling checkNewStyleAPI() means Login() resolves to
	// /api/auth/login without depending on how the console answers a GET of /.
	u.new = true

	client := &UNASClient{Unifi: u}

	if err := u.Login(); err != nil {
		return client, err
	}

	return client, nil
}

// GetUNASDeviceInfo returns console identity, CPU, and memory data.
func (c *UNASClient) GetUNASDeviceInfo() (*UNASDeviceInfo, error) {
	if c == nil || c.Unifi == nil {
		return nil, ErrNilUnifi
	}

	c.DebugLog("Polling UNAS for device info")

	response := &UNASDeviceInfo{}

	if err := c.GetData(APIUNASDeviceInfoPath, response); err != nil {
		return nil, fmt.Errorf("fetching UNAS device info: %w", err)
	}

	return response, nil
}

// GetUNASStorage returns storage pool and physical disk data.
func (c *UNASClient) GetUNASStorage() (*UNASStorage, error) {
	if c == nil || c.Unifi == nil {
		return nil, ErrNilUnifi
	}

	c.DebugLog("Polling UNAS for storage")

	response := &UNASStorage{}

	if err := c.GetData(APIUNASStoragePath, response); err != nil {
		return nil, fmt.Errorf("fetching UNAS storage: %w", err)
	}

	return response, nil
}

// GetUNASDrives returns the shares ("drives") configured on the console.
func (c *UNASClient) GetUNASDrives() ([]*UNASDrive, error) {
	if c == nil || c.Unifi == nil {
		return nil, ErrNilUnifi
	}

	c.DebugLog("Polling UNAS for drives")

	var response struct {
		Drives []UNASDrive `json:"drives"`
	}

	if err := c.GetData(APIUNASDrivesPath, &response); err != nil {
		return nil, fmt.Errorf("fetching UNAS drives: %w", err)
	}

	drives := make([]*UNASDrive, len(response.Drives))

	for i := range response.Drives {
		drives[i] = &response.Drives[i]
	}

	return drives, nil
}

// GetUNASNetworkIO returns instantaneous console network throughput.
func (c *UNASClient) GetUNASNetworkIO() (*UNASNetworkIO, error) {
	if c == nil || c.Unifi == nil {
		return nil, ErrNilUnifi
	}

	c.DebugLog("Polling UNAS for network IO")

	response := &UNASNetworkIO{}

	if err := c.GetData(APIUNASNetworkIOPath, response); err != nil {
		return nil, fmt.Errorf("fetching UNAS network IO: %w", err)
	}

	return response, nil
}

// GetUNASDevice polls every UNAS endpoint and returns the aggregate.
//
// Every endpoint is attempted even if an earlier one failed, so a console that exposes only
// some of them still reports what it has. A partial failure is deliberately NOT an error:
// the failing endpoints are logged to ErrorLog and their fields left nil, because the
// natural `if err != nil { return }` at the call site would otherwise drop every metric
// from a console whenever one endpoint 404s. An error is returned only when nothing at all
// could be fetched, and then the device is nil.
func (c *UNASClient) GetUNASDevice() (*UNASDevice, error) {
	if c == nil || c.Unifi == nil {
		return nil, ErrNilUnifi
	}

	device := &UNASDevice{SourceName: c.URL}
	errs := []error{}

	info, err := c.GetUNASDeviceInfo()
	if err != nil {
		errs = append(errs, err)
	} else {
		device.DeviceInfo = info
	}

	storage, err := c.GetUNASStorage()
	if err != nil {
		errs = append(errs, err)
	} else {
		device.Storage = storage
	}

	drives, err := c.GetUNASDrives()
	if err != nil {
		errs = append(errs, err)
	} else {
		device.Drives = drives
	}

	netIO, err := c.GetUNASNetworkIO()
	if err != nil {
		errs = append(errs, err)
	} else {
		device.NetworkIO = netIO
	}

	if len(errs) == unasEndpointCount {
		return nil, errors.Join(errs...)
	}

	for _, err := range errs {
		c.ErrorLog("UNAS %s: %v", c.URL, err)
	}

	return device, nil
}
