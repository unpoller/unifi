package unifi // nolint: testpackage

import (
	"encoding/pem"
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
	a.InDelta(87, s1.BatteryStatus.Percentage.Float64(), 0.001)
	require.NotNil(t, s1.Stats.Temperature)
	a.InDelta(21.5, s1.Stats.Temperature.Value.Float64(), 0.001)
	a.Equal("bridge-1", s1.WirelessConnectionState.Bridge)

	// A sensor with no battery reports a JSON null, which FlexFloat collapses to 0 -- this
	// is the accepted tradeoff of using Flex* types uniformly instead of nullable pointers.
	s2 := sensors[1]
	a.Equal("sensor-2", s2.ID)
	a.InDelta(0, s2.BatteryStatus.Percentage.Float64(), 0.001)
	a.Nil(s2.Stats.Temperature)
	a.Equal(int64(1734000500000), s2.LeakDetectedAt.Int64())
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
	a.True(cameras[0].HasPackageCamera.Val)
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
	a.True(lights[0].IsLightOn.Val)
	require.NotNil(t, lights[0].LightDeviceSettings)
	a.Equal(3, lights[0].LightDeviceSettings.LEDLevel.Int())
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
	a.Equal(32, bridges[0].MaxClients.Int())
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
	a.True(hub.IsAlarmHub.Val)
	require.NotNil(t, hub.AlarmHub)
	a.True(hub.AlarmHub.Armed.Val)
	require.NotNil(t, hub.AlarmHub.Battery)
	a.InDelta(13.2, hub.AlarmHub.Battery.Voltage.Float64(), 0.001)
	require.NotNil(t, hub.AlarmHub.InputPower)
	a.InDelta(12.1, hub.AlarmHub.InputPower.BT.Float64(), 0.001)
	// Typ1 is null in the fixture; FlexFloat collapses that to 0, same as a real 0 reading.
	a.InDelta(0, hub.AlarmHub.InputPower.Typ1.Float64(), 0.001)

	notHub := stations[1]
	a.False(notHub.IsAlarmHub.Val)
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
	a.Equal(0, nvr.ArmMode.BreachEventCount.Int())
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

// newProtectOnlyConsole returns a server that behaves like a UNVR: the Protect Integration
// paths answer from fixtures, /api/auth/login succeeds, and every Network path -- including
// the /proxy/network/status that NewUnifi ends on -- serves the UniFi OS SPA HTML that
// caused unpoller/unpoller#1066. The returned slice records requested paths, in order.
func newProtectOnlyConsole(t *testing.T) (*httptest.Server, *[]string) {
	t.Helper()

	requested := &[]string{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requested = append(*requested, r.URL.Path)

		if r.URL.Path == APILoginPathNew {
			w.Header().Set("x-csrf-token", "csrf-abc")
			w.WriteHeader(http.StatusOK)

			return
		}

		file, ok := protectFixtures[r.URL.Path]
		if !ok {
			// No Network application to route to, so UniFi OS serves its own SPA.
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<!doctype html><html lang="en"><body></body></html>`))

			return
		}

		body, err := os.ReadFile(file)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))

	t.Cleanup(srv.Close)

	return srv, requested
}

// TestNewUnifiFailsOnProtectOnlyConsole is the bug NewProtectClient exists to fix, pinned as
// a test: NewUnifi cannot construct a client for a console with no Network application. If
// this ever starts passing, NewProtectClient's reason for existing has changed.
func TestNewUnifiFailsOnProtectOnlyConsole(t *testing.T) {
	t.Parallel()

	srv, _ := newProtectOnlyConsole(t)

	_, err := NewUnifi(&Config{
		URL: srv.URL, User: "ro-user", Pass: "secret",
		ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to get server version")
}

func TestNewProtectClient(t *testing.T) {
	t.Parallel()

	srv, requested := newProtectOnlyConsole(t)

	c, err := NewProtectClient(&Config{
		URL: srv.URL, User: "ro-user", Pass: "secret",
		ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.NoError(t, err)
	require.NotNil(t, c)

	a := assert.New(t)
	// New-style paths are assumed rather than probed, so nothing but the login and the
	// reachability probe reaches the console -- and /status is never requested.
	a.True(c.new)
	a.Equal([]string{APILoginPathNew, APIProtectMetaInfoPath}, *requested)
	a.Equal("csrf-abc", c.csrf)

	// ServerStatus must be populated: Unifi embeds *ServerStatus, so a nil one turns every
	// caller's c.ServerVersion into a panic. The version is Protect's, not Network's.
	require.NotNil(t, c.ServerStatus)
	a.Equal("5.1.113", c.ServerVersion)
	a.True(c.Up.Val)
}

// TestNewProtectClientSkipsLoginWithAPIKey covers the one case where logging in is actively
// harmful: Login() routes through /status when Config.APIKey is set, and /status is exactly
// what a Protect-only console cannot answer.
func TestNewProtectClientSkipsLoginWithAPIKey(t *testing.T) {
	t.Parallel()

	srv, requested := newProtectOnlyConsole(t)

	c, err := NewProtectClient(&Config{
		URL: srv.URL, APIKey: "network-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.NoError(t, err)
	require.NotNil(t, c)

	assert.Equal(t, []string{APIProtectMetaInfoPath}, *requested)
}

// TestNewProtectClientSkipsLoginWithoutUser covers a key-only config: an operator with a
// Protect Integration key but no local account still gets a usable client.
func TestNewProtectClientSkipsLoginWithoutUser(t *testing.T) {
	t.Parallel()

	srv, requested := newProtectOnlyConsole(t)

	_, err := NewProtectClient(&Config{
		URL: srv.URL, ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{APIProtectMetaInfoPath}, *requested)
}

func TestNewProtectClientRejectsBadConfig(t *testing.T) {
	t.Parallel()

	_, err := NewProtectClient(nil)
	require.ErrorIs(t, err, ErrNilConfig)

	// Integration/v1 is X-API-Key only; user/pass alone can never authenticate it.
	_, err = NewProtectClient(&Config{URL: "https://127.0.0.1", User: "ro-user", Pass: "secret"})
	require.ErrorIs(t, err, ErrAPIKeyRequired)
}

// TestNewProtectClientProtectNotInstalled covers a console that answers neither application.
// The client is still returned so callers can inspect it, matching NewUnifi's contract.
func TestNewProtectClientProtectNotInstalled(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	c, err := NewProtectClient(&Config{
		URL: srv.URL, ProtectAPIKey: "protect-key-abc", DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.ErrorIs(t, err, ErrEndpointNotFound)
	require.NotNil(t, c)
	assert.Nil(t, c.ServerStatus)
}

// TestNewProtectClientPinsSSLCert covers the fingerprint loop NewProtectClient duplicates
// from NewUnifi. It writes into a slice newUnifi only allocates when SSLCert is non-empty,
// so an empty-cert regression there is an index panic rather than a failure elsewhere.
func TestNewProtectClientPinsSSLCert(t *testing.T) {
	t.Parallel()

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"applicationVersion":"7.2.105"}`))
	}))
	t.Cleanup(srv.Close)

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})
	require.NotNil(t, cert)

	c, err := NewProtectClient(&Config{
		URL: srv.URL, ProtectAPIKey: "protect-key-abc", SSLCert: [][]byte{cert},
		DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.NoError(t, err)
	require.NotNil(t, c)

	a := assert.New(t)
	a.Len(c.fingerprints, 1)
	a.NotNil(c.Transport)
	a.Equal("7.2.105", c.ServerVersion)
}
