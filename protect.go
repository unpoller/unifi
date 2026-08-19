package unifi

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Protect device support is built on Ubiquiti's official, documented Protect Integration
// API (https://developer.ui.com/protect/), not the private bootstrap API. See
// unpoller/unpoller#1015. Getters are deliberately not part of the UnifiClient interface,
// matching the UNAS convention in unas.go: a Network-only console answers none of them, and
// adding them to the interface would force matching methods onto mocks.MockUnifi for no
// benefit, since inputunifi holds a concrete *Unifi rather than the interface.

// protectEndpointCount is the number of endpoints GetProtectDevices polls: sensors, cameras,
// lights, bridges, link-stations, and the NVR. /v1/meta/info is a separate reachability
// probe (used by callers to detect "Protect not installed" before polling) and is not part
// of this tally.
const protectEndpointCount = 6

// ProtectStateValue maps a Protect device connection state to a numeric gauge value, so every
// output plugin agrees on the mapping instead of each inventing its own. Keep the raw state
// string as a label/tag alongside the number: the number is for alerting, the string for
// debugging. Unknown states are -1, never 0 -- 0 means DISCONNECTED.
func ProtectStateValue(state string) float64 {
	switch state {
	case "CONNECTED":
		return 2
	case "CONNECTING":
		return 1
	case "DISCONNECTED":
		return 0
	default:
		return -1
	}
}

// ProtectMetaInfo holds the Protect application version, from /v1/meta/info. It is the
// cheapest available endpoint and the recommended probe for "is Protect installed and is
// this key valid" before polling the rest.
type ProtectMetaInfo struct {
	ApplicationVersion string `json:"applicationVersion"`
}

// ProtectDeviceIdentity holds the fields common to every Protect device.
type ProtectDeviceIdentity struct {
	ID       string `json:"id"`
	ModelKey string `json:"modelKey"`
	State    string `json:"state"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	GUID     string `json:"guid"`
	MAC      string `json:"mac"`
}

// ProtectBatteryStatus holds a device's battery percentage and low-battery flag.
// Percentage is nullable in the spec: a device with no battery reports Percentage == nil,
// which must be skipped by callers, never treated as 0.
type ProtectBatteryStatus struct {
	Percentage *float64 `json:"percentage"`
	IsLow      bool     `json:"isLow"`
}

// ProtectSignalState holds a wireless device's signal quality and strength.
type ProtectSignalState struct {
	SignalQuality  *float64 `json:"signalQuality"`
	SignalStrength *float64 `json:"signalStrength"`
}

// ProtectWirelessConnectionState holds a wireless sensor's link to its bridge.
type ProtectWirelessConnectionState struct {
	SignalState   *ProtectSignalState   `json:"signalState"`
	BatteryStatus *ProtectBatteryStatus `json:"batteryStatus"`
	Bridge        string                `json:"bridge"`
}

// ProtectSensorChannelFlag reports whether a sensor supports a given measurement and, if so,
// on how many channels.
type ProtectSensorChannelFlag struct {
	ChannelCount int `json:"channelCount"`
}

// ProtectSensorFeatureFlags reports which measurements a sensor supports.
type ProtectSensorFeatureFlags struct {
	Temperature *ProtectSensorChannelFlag `json:"temperature"`
	Humidity    *ProtectSensorChannelFlag `json:"humidity"`
	Light       *ProtectSensorChannelFlag `json:"light"`
	Motion      *ProtectSensorChannelFlag `json:"motion"`
	WaterLeak   *ProtectSensorChannelFlag `json:"waterLeak"`
	Open        *ProtectSensorChannelFlag `json:"open"`
	Tamper      *ProtectSensorChannelFlag `json:"tamper"`
	Smoke       *ProtectSensorChannelFlag `json:"smoke"`
}

// ProtectSensorStatValue holds one sensor reading. Value is nullable: a sensor without this
// feature, or one that hasn't reported yet, sends null and must be skipped, not zeroed.
type ProtectSensorStatValue struct {
	Value  *float64 `json:"value"`
	Status string   `json:"status"` // neutral, low, safe, high, unknown
}

// ProtectSensorStats holds a sensor's current readings.
type ProtectSensorStats struct {
	Light       *ProtectSensorStatValue `json:"light"`
	Humidity    *ProtectSensorStatValue `json:"humidity"`
	Temperature *ProtectSensorStatValue `json:"temperature"`
}

// ProtectSensorThresholdSettings holds the alerting configuration for one sensor measurement.
type ProtectSensorThresholdSettings struct {
	IsEnabled     bool     `json:"isEnabled"`
	Margin        *float64 `json:"margin"`
	LowThreshold  *float64 `json:"lowThreshold"`
	HighThreshold *float64 `json:"highThreshold"`
}

// ProtectSensorMotionSettings holds a sensor's motion-detection configuration.
type ProtectSensorMotionSettings struct {
	IsEnabled            bool `json:"isEnabled"`
	Sensitivity          int  `json:"sensitivity"`
	SensitivityWhenArmed int  `json:"sensitivityWhenArmed"`
}

// ProtectSensor represents a Protect sensor, including SuperLink sensors (unpoller/unpoller#1015).
// API Path: /proxy/protect/integration/v1/sensors
type ProtectSensor struct {
	ProtectDeviceIdentity

	MountType string `json:"mountType"` // door, window, garage, leak, none

	BatteryStatus *ProtectBatteryStatus      `json:"batteryStatus"`
	FeatureFlags  *ProtectSensorFeatureFlags `json:"featureFlags"`
	Stats         *ProtectSensorStats        `json:"stats"`

	LightSettings       *ProtectSensorThresholdSettings `json:"lightSettings"`
	HumiditySettings    *ProtectSensorThresholdSettings `json:"humiditySettings"`
	TemperatureSettings *ProtectSensorThresholdSettings `json:"temperatureSettings"`

	IsOpened            bool   `json:"isOpened"`
	OpenStatusChangedAt *int64 `json:"openStatusChangedAt"`

	IsMotionDetected bool                         `json:"isMotionDetected"`
	MotionDetectedAt *int64                       `json:"motionDetectedAt"`
	MotionSettings   *ProtectSensorMotionSettings `json:"motionSettings"`

	ScheduleMode  string   `json:"scheduleMode"` // always, when_armed
	ArmProfileIDs []string `json:"armProfileIds"`

	AlarmTriggeredAt       *int64 `json:"alarmTriggeredAt"`
	LeakDetectedAt         *int64 `json:"leakDetectedAt"`
	ExternalLeakDetectedAt *int64 `json:"externalLeakDetectedAt"`
	TamperingDetectedAt    *int64 `json:"tamperingDetectedAt"`

	WirelessConnectionState *ProtectWirelessConnectionState `json:"wirelessConnectionState"`

	// GlassBreakSettings has no documented shape; kept raw for diagnosis, per the UNAS
	// convention (unas.go).
	GlassBreakSettings json.RawMessage `fake:"skip" json:"glassBreakSettings"`
}

// ProtectCamera represents a Protect camera. Only state is metric-shaped in v1; the rest is
// configuration, per the scope decision in the Protect device metrics plan.
// API Path: /proxy/protect/integration/v1/cameras
type ProtectCamera struct {
	ProtectDeviceIdentity

	IsMicEnabled     bool   `json:"isMicEnabled"`
	MicVolume        int    `json:"micVolume"`
	ActivePatrolSlot *int   `json:"activePatrolSlot"`
	VideoMode        string `json:"videoMode"`
	HdrType          string `json:"hdrType"`
	HasPackageCamera bool   `json:"hasPackageCamera"`

	// The following are documented as objects but their detailed shape is not; kept raw for
	// diagnosis rather than invented, per the UNAS convention (unas.go).
	OSDSettings         json.RawMessage `fake:"skip" json:"osdSettings"`
	LEDSettings         json.RawMessage `fake:"skip" json:"ledSettings"`
	LCDMessage          json.RawMessage `fake:"skip" json:"lcdMessage"`
	FeatureFlags        json.RawMessage `fake:"skip" json:"featureFlags"`
	SmartDetectSettings json.RawMessage `fake:"skip" json:"smartDetectSettings"`
}

// ProtectLightDeviceSettings holds a light's physical configuration.
type ProtectLightDeviceSettings struct {
	IsIndicatorEnabled bool `json:"isIndicatorEnabled"`
	PIRDuration        int  `json:"pirDuration"`
	PIRSensitivity     int  `json:"pirSensitivity"`
	LEDLevel           int  `json:"ledLevel"`
}

// ProtectLight represents a Protect floodlight.
// API Path: /proxy/protect/integration/v1/lights
type ProtectLight struct {
	ProtectDeviceIdentity

	IsDark              bool   `json:"isDark"`
	IsLightOn           bool   `json:"isLightOn"`
	IsLightForceEnabled bool   `json:"isLightForceEnabled"`
	IsPirMotionDetected bool   `json:"isPirMotionDetected"`
	LastMotion          *int64 `json:"lastMotion"`
	Camera              string `json:"camera"` // ID of the paired camera, if any

	LightDeviceSettings *ProtectLightDeviceSettings `json:"lightDeviceSettings"`

	// LightModeSettings is documented as an object but its detailed shape is not; kept raw
	// for diagnosis rather than invented, per the UNAS convention (unas.go).
	LightModeSettings json.RawMessage `fake:"skip" json:"lightModeSettings"`
}

// ProtectBridge represents a Protect wireless bridge, such as a SuperLink Gateway.
// API Path: /proxy/protect/integration/v1/bridges
type ProtectBridge struct {
	ProtectDeviceIdentity

	Platform   string   `json:"platform"`
	Clients    []string `json:"clients"`
	MaxClients int      `json:"maxClients"`
}

// ProtectLinkStationCover holds a link station's tamper-cover state.
type ProtectLinkStationCover struct {
	Distance *float64 `json:"distance"`
	Status   string   `json:"status"`
}

// ProtectLinkStationBattery holds a link station's (alarm hub's) battery state.
type ProtectLinkStationBattery struct {
	Charging      bool     `json:"charging"`
	Connection    string   `json:"connection"`
	Voltage       *float64 `json:"voltage"`
	BatteryStatus string   `json:"batteryStatus"` // ok, low, critical
}

// ProtectInputPower holds a link station's input power readings, by input type.
type ProtectInputPower struct {
	BT   *float64 `json:"bt"`
	Typ1 *float64 `json:"typ1"`
	Typ2 *float64 `json:"typ2"`
}

// ProtectAlarmHub holds the metrics exposed by a link station acting as an alarm hub
// (isAlarmHub == true), such as a SuperLink Gateway.
type ProtectAlarmHub struct {
	Armed      bool                       `json:"armed"`
	Battery    *ProtectLinkStationBattery `json:"battery"`
	Cover      *ProtectLinkStationCover   `json:"cover"`
	InputPower *ProtectInputPower         `json:"inputPower"`

	// The following are documented as objects but their detailed shape is not; kept raw for
	// diagnosis rather than invented, per the UNAS convention (unas.go).
	BuckBoost                    json.RawMessage `fake:"skip" json:"buckboost"`
	Connector                    json.RawMessage `fake:"skip" json:"connector"`
	CurrentMeterChannelStatus    json.RawMessage `fake:"skip" json:"currentMeterChannelStatus"`
	CurrentMeterStatus           json.RawMessage `fake:"skip" json:"currentMeterStatus"`
	Poeout                       json.RawMessage `fake:"skip" json:"poeout"`
	PowerMeter                   json.RawMessage `fake:"skip" json:"powerMeter"`
	Output                       json.RawMessage `fake:"skip" json:"output"`
	Input                        json.RawMessage `fake:"skip" json:"input"`
	InputTerminalStatus          json.RawMessage `fake:"skip" json:"inputTerminalStatus"`
	OutputTerminalStatus         json.RawMessage `fake:"skip" json:"outputTerminalStatus"`
	EmergencyTerminalStatus      json.RawMessage `fake:"skip" json:"emergencyTerminalStatus"`
	AuxiliaryPowerTerminalStatus json.RawMessage `fake:"skip" json:"auxiliaryPowerTerminalStatus"`
}

// ProtectLinkStation represents a Protect link station. /v1/alarm-hubs returns the same
// schema filtered to IsAlarmHub == true; GetProtectLinkStations covers both and the
// alarm-hubs endpoint is deliberately not polled, to avoid double-counting.
// API Path: /proxy/protect/integration/v1/link-stations
type ProtectLinkStation struct {
	ProtectDeviceIdentity

	IsAlarmHub bool             `json:"isAlarmHub"`
	LastEvent  *int64           `json:"lastEvent"`
	AlarmHub   *ProtectAlarmHub `json:"alarmHub"`

	// LedSettings is documented as an object but its detailed shape is not; kept raw for
	// diagnosis rather than invented, per the UNAS convention (unas.go).
	LedSettings json.RawMessage `fake:"skip" json:"ledSettings"`
}

// ProtectArmMode holds an NVR's alarm-arming state.
type ProtectArmMode struct {
	Status               string `json:"status"` // arming, armed, breach, disabled
	ArmProfileID         string `json:"armProfileId"`
	ArmedAt              *int64 `json:"armedAt"`
	WillBeArmedAt        *int64 `json:"willBeArmedAt"`
	BreachDetectedAt     *int64 `json:"breachDetectedAt"`
	BreachEventCount     int    `json:"breachEventCount"`
	BreachTriggerEventID string `json:"breachTriggerEventId"`
	BreachEventID        string `json:"breachEventId"`
}

// ProtectNVR represents the NVR console itself. Only ArmMode is metric-shaped in v1; the
// rest is configuration. GET /v1/nvrs returns a single object, not an array.
// API Path: /proxy/protect/integration/v1/nvrs
type ProtectNVR struct {
	ProtectDeviceIdentity

	ArmMode *ProtectArmMode `json:"armMode"`

	// DoorbellSettings is documented as an object but its detailed shape is not; kept raw
	// for diagnosis rather than invented, per the UNAS convention (unas.go).
	DoorbellSettings json.RawMessage `fake:"skip" json:"doorbellSettings"`
}

// ProtectDevices is the aggregate of every polled Protect endpoint for one console.
// Any member may be nil (or empty, for slices) if that endpoint failed; check before use.
type ProtectDevices struct {
	SourceName string
	// Version is not populated by GetProtectDevices -- adding a meta/info call here would
	// make it a 7th polled endpoint and break protectEndpointCount accounting. Callers that
	// already probe /v1/meta/info for reachability (see ProtectMetaInfo) should set this
	// themselves from that result.
	Version string
	NVR     *ProtectNVR

	Sensors      []*ProtectSensor      `fakesize:"5"`
	Cameras      []*ProtectCamera      `fakesize:"5"`
	Lights       []*ProtectLight       `fakesize:"5"`
	Bridges      []*ProtectBridge      `fakesize:"5"`
	LinkStations []*ProtectLinkStation `fakesize:"5"`
}

// GetProtectMetaInfo returns the Protect application version. This is the cheapest
// available endpoint and the recommended probe for "is Protect installed and is this key
// valid" before polling the rest.
func (u *Unifi) GetProtectMetaInfo() (*ProtectMetaInfo, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for meta info")

	info := &ProtectMetaInfo{}

	if err := u.GetData(APIProtectMetaInfoPath, info); err != nil {
		return nil, fmt.Errorf("fetching Protect meta info: %w", err)
	}

	return info, nil
}

// getProtectList fetches a Protect Integration/v1 list endpoint. Unlike the Network
// Integration/v1 API (see getIntegrationList in integration.go), Protect list endpoints
// return a bare JSON array with no pagination envelope.
func getProtectList[T any](u *Unifi, path string) ([]*T, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	var items []T

	if err := u.GetData(path, &items); err != nil {
		return nil, err
	}

	result := make([]*T, len(items))

	for i := range items {
		result[i] = &items[i]
	}

	return result, nil
}

// GetProtectSensors returns all Protect sensors, including SuperLink sensors.
func (u *Unifi) GetProtectSensors() ([]*ProtectSensor, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for sensors")

	sensors, err := getProtectList[ProtectSensor](u, APIProtectSensorsPath)
	if err != nil {
		return nil, fmt.Errorf("fetching Protect sensors: %w", err)
	}

	return sensors, nil
}

// GetProtectCameras returns all Protect cameras.
func (u *Unifi) GetProtectCameras() ([]*ProtectCamera, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for cameras")

	cameras, err := getProtectList[ProtectCamera](u, APIProtectCamerasPath)
	if err != nil {
		return nil, fmt.Errorf("fetching Protect cameras: %w", err)
	}

	return cameras, nil
}

// GetProtectLights returns all Protect lights.
func (u *Unifi) GetProtectLights() ([]*ProtectLight, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for lights")

	lights, err := getProtectList[ProtectLight](u, APIProtectLightsPath)
	if err != nil {
		return nil, fmt.Errorf("fetching Protect lights: %w", err)
	}

	return lights, nil
}

// GetProtectBridges returns all Protect wireless bridges, such as SuperLink Gateways.
func (u *Unifi) GetProtectBridges() ([]*ProtectBridge, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for bridges")

	bridges, err := getProtectList[ProtectBridge](u, APIProtectBridgesPath)
	if err != nil {
		return nil, fmt.Errorf("fetching Protect bridges: %w", err)
	}

	return bridges, nil
}

// GetProtectLinkStations returns all Protect link stations, including alarm hubs.
// /v1/alarm-hubs is deliberately not polled: it returns the same schema filtered to
// IsAlarmHub == true, and polling both would double-count every alarm hub.
func (u *Unifi) GetProtectLinkStations() ([]*ProtectLinkStation, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for link stations")

	stations, err := getProtectList[ProtectLinkStation](u, APIProtectLinkStationsPath)
	if err != nil {
		return nil, fmt.Errorf("fetching Protect link stations: %w", err)
	}

	return stations, nil
}

// GetProtectNVR returns the NVR console itself. GET /v1/nvrs returns a single object, not
// an array -- easy to get wrong.
func (u *Unifi) GetProtectNVR() (*ProtectNVR, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	u.DebugLog("Polling Protect Integration/v1 for NVR")

	nvr := &ProtectNVR{}

	if err := u.GetData(APIProtectNVRPath, nvr); err != nil {
		return nil, fmt.Errorf("fetching Protect NVR: %w", err)
	}

	return nvr, nil
}

// GetProtectDevices polls every Protect endpoint and returns the aggregate.
//
// Every endpoint is attempted even if an earlier one failed, so a console that exposes only
// some of them still reports what it has. A partial failure is deliberately NOT an error:
// the failing endpoints are logged to ErrorLog and their fields left nil/empty, because the
// natural `if err != nil { return }` at the call site would otherwise drop every metric
// whenever one endpoint 404s -- and a Protect estate with no lights or no link-stations is
// the normal case, not an error. An error is returned only when nothing at all could be
// fetched, and then the device is nil. This mirrors GetUNASDevice in unas.go.
func (u *Unifi) GetProtectDevices() (*ProtectDevices, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	device := &ProtectDevices{SourceName: u.URL}
	errs := []error{}

	if sensors, err := u.GetProtectSensors(); err != nil {
		errs = append(errs, err)
	} else {
		device.Sensors = sensors
	}

	if cameras, err := u.GetProtectCameras(); err != nil {
		errs = append(errs, err)
	} else {
		device.Cameras = cameras
	}

	if lights, err := u.GetProtectLights(); err != nil {
		errs = append(errs, err)
	} else {
		device.Lights = lights
	}

	if bridges, err := u.GetProtectBridges(); err != nil {
		errs = append(errs, err)
	} else {
		device.Bridges = bridges
	}

	if stations, err := u.GetProtectLinkStations(); err != nil {
		errs = append(errs, err)
	} else {
		device.LinkStations = stations
	}

	if nvr, err := u.GetProtectNVR(); err != nil {
		errs = append(errs, err)
	} else {
		device.NVR = nvr
	}

	if len(errs) == protectEndpointCount {
		return nil, errors.Join(errs...)
	}

	for _, err := range errs {
		u.ErrorLog("Protect %s: %v", u.URL, err)
	}

	return device, nil
}
