package unifi // nolint: testpackage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// protectFixtures maps each Protect API path to the file serving its response.
var protectFixtures = map[string]string{
	APIProtectMetaInfoPath:     "endpoints_data/protect-meta-info.json",
	APIProtectSensorsPath:      "endpoints_data/protect-sensors.json",
	APIProtectCamerasPath:      "endpoints_data/protect-cameras.json",
	APIProtectLightsPath:       "endpoints_data/protect-lights.json",
	APIProtectBridgesPath:      "endpoints_data/protect-bridges.json",
	APIProtectLinkStationsPath: "endpoints_data/protect-link-stations.json",
	APIProtectNVRPath:          "endpoints_data/protect-nvrs.json",
}

// newProtectTestClient returns a *Unifi pointed at a server that answers the Protect paths
// from fixtures. Paths listed in fail are answered with a 500 instead. The returned slice
// records the paths and X-API-Key headers the client actually requested, in order.
func newProtectTestClient(t *testing.T, fail ...string) (*Unifi, *[]string) {
	t.Helper()

	failed := make(map[string]bool, len(fail))
	for _, p := range fail {
		failed[p] = true
	}

	requested := &[]string{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requested = append(*requested, r.URL.Path)

		if failed[r.URL.Path] {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		file, ok := protectFixtures[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		body, err := os.ReadFile(file)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))

	t.Cleanup(srv.Close)

	return &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
	}, requested
}

func TestProtectGetMetaInfo(t *testing.T) {
	t.Parallel()

	c, requested := newProtectTestClient(t)

	info, err := c.GetProtectMetaInfo()
	require.NoError(t, err)
	require.NotNil(t, info)

	a := assert.New(t)
	// The /proxy/ prefix must survive path(): Protect has no /proxy/network app.
	a.Equal([]string{APIProtectMetaInfoPath}, *requested)
	a.Equal("5.1.113", info.ApplicationVersion)
}

// TestProtectMetaInfoNotInstalled confirms a 404 (the "Protect not installed" case) still
// surfaces as ErrEndpointNotFound after GetProtectMetaInfo's fmt.Errorf wrap. Part 2's
// collectProtect depends on errors.Is matching through that wrap to detect this case cleanly.
func TestProtectMetaInfoNotInstalled(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	c := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
	}

	_, err := c.GetProtectMetaInfo()
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEndpointNotFound)
}

func TestProtectAuthUsesProtectAPIKey(t *testing.T) {
	t.Parallel()

	var gotAPIKey, gotCSRF string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("X-API-Key")
		gotCSRF = r.Header.Get("X-CSRF-Token")

		body, err := os.ReadFile("endpoints_data/protect-meta-info.json")
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
		csrf:   "network-csrf-should-not-be-used",
	}

	_, err := c.GetProtectMetaInfo()
	require.NoError(t, err)

	a := assert.New(t)
	a.Equal("protect-key-abc", gotAPIKey)
	a.Empty(gotCSRF)
}

func TestProtectAuthFallsBackToAPIKey(t *testing.T) {
	t.Parallel()

	var gotAPIKey string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("X-API-Key")

		body, err := os.ReadFile("endpoints_data/protect-meta-info.json")
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, APIKey: "network-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
	}

	_, err := c.GetProtectMetaInfo()
	require.NoError(t, err)
	assert.Equal(t, "network-key-abc", gotAPIKey)
}

func TestProtectGetSensors(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t)

	sensors, err := c.GetProtectSensors()
	require.NoError(t, err)
	require.Len(t, sensors, 2)

	a := assert.New(t)

	s1 := sensors[0]
	a.Equal("sensor-1", s1.ID)
	a.Equal("sensor", s1.ModelKey)
	a.Equal("CONNECTED", s1.State)
	a.Equal("door", s1.MountType)
	require.NotNil(t, s1.BatteryStatus.Percentage)
	a.InDelta(87, *s1.BatteryStatus.Percentage, 0.001)
	require.NotNil(t, s1.Stats.Temperature.Value)
	a.InDelta(21.5, *s1.Stats.Temperature.Value, 0.001)
	a.Equal("bridge-1", s1.WirelessConnectionState.Bridge)

	// A sensor with no battery reports a nil pointer, never a zeroed 0 -- that distinction
	// is the whole reason Protect fields use pointers instead of FlexInt.
	s2 := sensors[1]
	a.Equal("sensor-2", s2.ID)
	a.Nil(s2.BatteryStatus.Percentage)
	a.Nil(s2.Stats.Temperature)
	require.NotNil(t, s2.LeakDetectedAt)
	a.Equal(int64(1734000500000), *s2.LeakDetectedAt)
}

func TestProtectGetCameras(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t)

	cameras, err := c.GetProtectCameras()
	require.NoError(t, err)
	require.Len(t, cameras, 2)

	a := assert.New(t)
	a.Equal("camera-1", cameras[0].ID)
	a.Equal("CONNECTED", cameras[0].State)
	a.True(cameras[0].HasPackageCamera)
	a.Equal("DISCONNECTED", cameras[1].State)
}

func TestProtectGetLights(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t)

	lights, err := c.GetProtectLights()
	require.NoError(t, err)
	require.Len(t, lights, 1)

	a := assert.New(t)
	a.Equal("light-1", lights[0].ID)
	a.True(lights[0].IsLightOn)
	require.NotNil(t, lights[0].LightDeviceSettings)
	a.Equal(3, lights[0].LightDeviceSettings.LEDLevel)
}

func TestProtectGetBridges(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t)

	bridges, err := c.GetProtectBridges()
	require.NoError(t, err)
	require.Len(t, bridges, 1)

	a := assert.New(t)
	a.Equal("bridge-1", bridges[0].ID)
	a.Len(bridges[0].Clients, 2)
	a.Equal(32, bridges[0].MaxClients)
}

func TestProtectGetLinkStations(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t)

	stations, err := c.GetProtectLinkStations()
	require.NoError(t, err)
	require.Len(t, stations, 2)

	a := assert.New(t)

	hub := stations[0]
	a.Equal("linkstation-1", hub.ID)
	a.True(hub.IsAlarmHub)
	require.NotNil(t, hub.AlarmHub)
	a.True(hub.AlarmHub.Armed)
	require.NotNil(t, hub.AlarmHub.Battery.Voltage)
	a.InDelta(13.2, *hub.AlarmHub.Battery.Voltage, 0.001)
	require.NotNil(t, hub.AlarmHub.InputPower.BT)
	a.InDelta(12.1, *hub.AlarmHub.InputPower.BT, 0.001)
	a.Nil(hub.AlarmHub.InputPower.Typ1)

	notHub := stations[1]
	a.False(notHub.IsAlarmHub)
	a.Nil(notHub.AlarmHub)
}

func TestProtectGetNVR(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t)

	nvr, err := c.GetProtectNVR()
	require.NoError(t, err)
	require.NotNil(t, nvr)

	a := assert.New(t)
	a.Equal("nvr-1", nvr.ID)
	require.NotNil(t, nvr.ArmMode)
	a.Equal("armed", nvr.ArmMode.Status)
	a.Equal(0, nvr.ArmMode.BreachEventCount)
}

func TestProtectGetDevices(t *testing.T) {
	t.Parallel()

	c, requested := newProtectTestClient(t)

	device, err := c.GetProtectDevices()
	require.NoError(t, err)
	require.NotNil(t, device)

	a := assert.New(t)
	a.Len(*requested, protectEndpointCount)
	a.NotEmpty(device.SourceName)
	a.Len(device.Sensors, 2)
	a.Len(device.Cameras, 2)
	a.Len(device.Lights, 1)
	a.Len(device.Bridges, 1)
	a.Len(device.LinkStations, 2)
	a.NotNil(device.NVR)
}

func TestProtectGetDevicesPartialFailure(t *testing.T) {
	t.Parallel()

	c, requested := newProtectTestClient(t, APIProtectLightsPath)

	logged := []string{}
	c.ErrorLog = func(msg string, v ...any) {
		logged = append(logged, fmt.Sprintf(msg, v...))
	}

	device, err := c.GetProtectDevices()
	// A console that will not serve one endpoint still reports the rest, and a partial
	// failure is deliberately not an error -- an `if err != nil { return }` at the call
	// site would otherwise discard every metric this console can still supply.
	require.NoError(t, err)
	require.NotNil(t, device)

	a := assert.New(t)
	a.Len(logged, 1, "the endpoint that failed is reported through ErrorLog instead")
	a.Contains(logged[0], "lights")
	a.Len(*requested, protectEndpointCount, "every endpoint is attempted even after one fails")
	a.Len(device.Sensors, 2)
	a.Empty(device.Lights)
	a.NotNil(device.NVR)
}

func TestProtectGetDevicesTotalFailure(t *testing.T) {
	t.Parallel()

	c, _ := newProtectTestClient(t,
		APIProtectSensorsPath, APIProtectCamerasPath, APIProtectLightsPath,
		APIProtectBridgesPath, APIProtectLinkStationsPath, APIProtectNVRPath)

	device, err := c.GetProtectDevices()
	// Total failure is now the only case that returns an error, and it returns no device.
	require.Error(t, err)
	require.Nil(t, device)
	// The joined error names every endpoint, so an operator sees the whole picture at once.
	require.Contains(t, err.Error(), "sensors")
	require.Contains(t, err.Error(), "cameras")
	require.Contains(t, err.Error(), "lights")
	require.Contains(t, err.Error(), "bridges")
	require.Contains(t, err.Error(), "link stations")
	require.Contains(t, err.Error(), "NVR")
}

func TestProtectNilReceivers(t *testing.T) {
	t.Parallel()

	var c *Unifi

	a := assert.New(t)

	_, err := c.GetProtectDevices()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectMetaInfo()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectSensors()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectCameras()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectLights()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectBridges()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectLinkStations()
	a.ErrorIs(err, ErrNilUnifi)

	_, err = c.GetProtectNVR()
	a.ErrorIs(err, ErrNilUnifi)
}

func TestProtectStateValue(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	a.InDelta(2.0, ProtectStateValue("CONNECTED"), 0.001)
	a.InDelta(1.0, ProtectStateValue("CONNECTING"), 0.001)
	a.InDelta(0.0, ProtectStateValue("DISCONNECTED"), 0.001)
	a.InDelta(-1.0, ProtectStateValue("UNKNOWN_STATE"), 0.001)
	a.InDelta(-1.0, ProtectStateValue(""), 0.001)
}

func TestProtectStructs(t *testing.T) {
	t.Parallel()

	var sensor ProtectSensor

	require.NoError(t, gofakeit.Struct(&sensor))
	require.NotEmpty(t, sensor.ID)

	var camera ProtectCamera

	require.NoError(t, gofakeit.Struct(&camera))

	var light ProtectLight

	require.NoError(t, gofakeit.Struct(&light))

	var bridge ProtectBridge

	require.NoError(t, gofakeit.Struct(&bridge))

	var station ProtectLinkStation

	require.NoError(t, gofakeit.Struct(&station))

	var nvr ProtectNVR

	require.NoError(t, gofakeit.Struct(&nvr))

	var devices ProtectDevices

	require.NoError(t, gofakeit.Struct(&devices))
	require.Len(t, devices.Sensors, 5)
}
