package unifi // nolint: testpackage

import (
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unasFixtures maps each UNAS API path to the file serving its response.
var unasFixtures = map[string]string{
	APIUNASDeviceInfoPath: "endpoints_data/unas-device-info.json",
	APIUNASStoragePath:    "endpoints_data/unas-storage.json",
	APIUNASDrivesPath:     "endpoints_data/unas-drives.json",
	APIUNASNetworkIOPath:  "endpoints_data/unas-network-io.json",
}

// newUNASTestClient returns a client pointed at a server that answers the UNAS paths from
// fixtures. Paths listed in fail are answered with a 500 instead. The returned slice
// records the paths the client actually requested, in order.
func newUNASTestClient(t *testing.T, fail ...string) (*UNASClient, *[]string) {
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

		file, ok := unasFixtures[r.URL.Path]
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

	return &UNASClient{Unifi: &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
	}}, requested
}

func TestUNASGetDeviceInfo(t *testing.T) {
	t.Parallel()

	c, requested := newUNASTestClient(t)

	info, err := c.GetUNASDeviceInfo()
	require.NoError(t, err)
	require.NotNil(t, info)

	a := assert.New(t)
	// The /proxy/ prefix must survive path(): a UNAS console has no /proxy/network app.
	a.Equal([]string{APIUNASDeviceInfoPath}, *requested)
	a.Equal("UNAS Pro", info.Name)
	a.Equal("UNASPro", info.Model)
	a.Equal("4.3.5", info.Version)
	a.Equal("4.3.5.1234", info.FirmwareVersion)
	a.Equal("healthy", info.Status)
	a.False(info.SfpAggregation.Val)
	a.InDelta(4.75, info.CPU.CurrentLoad.Val, 0.001)
	a.Equal(46, info.CPU.Temperature.Int())
	a.Equal(int64(8589934592), info.Memory.Total.Int64())
	a.Equal(int64(3221225472), info.Memory.Free.Int64())
	a.Equal(int64(5368709120), info.Memory.Available.Int64())

	// The reference implementation shipped an empty json tag here, so this never parsed.
	require.Len(t, info.NetworkInterfaces, 2)
	a.Equal("eth0", info.NetworkInterfaces[0].Interface)
	a.True(info.NetworkInterfaces[0].Connected.Val)
	a.False(info.NetworkInterfaces[1].Connected.Val)
}

func TestUNASGetStorage(t *testing.T) {
	t.Parallel()

	c, _ := newUNASTestClient(t)

	storage, err := c.GetUNASStorage()
	require.NoError(t, err)
	require.NotNil(t, storage)
	require.Len(t, storage.Pools, 1)
	require.Len(t, storage.Disks, 2)

	a := assert.New(t)
	pool := storage.Pools[0]
	a.Equal("pool-1", pool.ID)
	a.Equal("healthy", pool.Status)
	a.Equal(int64(60000000000000), pool.Capacity.Int64())
	a.Equal(int64(12500000000000), pool.Usage.Int64())
	// The reference implementation left this field unexported, so it never parsed.
	a.Equal("rg-1", pool.ActiveRaidGroupID)
	require.Len(t, pool.RaidGroups, 1)
	a.Equal("RAID5", pool.RaidGroups[0].CurrentLevel)
	a.Equal(100, pool.RaidGroups[0].Progress.Int())

	disk := storage.Disks[1]
	a.Equal("2", disk.SlotID)
	a.Equal("pool-1", disk.PoolID)
	a.Equal("SERIAL0002", disk.Serial)
	a.Equal(7200, disk.RPM.Int())
	a.Equal(39, disk.Temperature.Int())
	a.Equal(8761, disk.PowerOnHours.Int())
	a.Equal(2, disk.BadSectorCount.Int())
	a.Equal(98, disk.HealthScore.Int())
	a.InDelta(130, disk.ReadKBPS.Val, 0.001)
	a.True(disk.SmartTestSupported.Val)
}

func TestUNASGetDrives(t *testing.T) {
	t.Parallel()

	c, _ := newUNASTestClient(t)

	drives, err := c.GetUNASDrives()
	require.NoError(t, err)
	require.Len(t, drives, 1)

	a := assert.New(t)
	a.Equal("drive-1", drives[0].ID)
	a.Equal("media", drives[0].Name)
	a.Equal("pool-1", drives[0].StoragePoolID)
	a.Equal(int64(10000000000000), drives[0].Quota.Int64())
	a.Equal(int64(4200000000000), drives[0].Usage.Int64())
	a.Equal(3, drives[0].MemberCount.Int())
	a.True(drives[0].Protections.SnapshotEnabled.Val)
	a.False(drives[0].Protections.RemoteBackupEnabled.Val)
}

func TestUNASGetNetworkIO(t *testing.T) {
	t.Parallel()

	c, _ := newUNASTestClient(t)

	netIO, err := c.GetUNASNetworkIO()
	require.NoError(t, err)
	require.NotNil(t, netIO)

	a := assert.New(t)
	a.InDelta(1024.5, netIO.ReceiveKBPS.Val, 0.001)
	a.InDelta(512.25, netIO.TransmitKBPS.Val, 0.001)
	a.Equal("2026-08-18T12:00:00Z", netIO.Timestamp)
}

func TestUNASGetDevice(t *testing.T) {
	t.Parallel()

	c, requested := newUNASTestClient(t)

	device, err := c.GetUNASDevice()
	require.NoError(t, err)
	require.NotNil(t, device)

	a := assert.New(t)
	a.Len(*requested, unasEndpointCount)
	a.Equal("UNAS Pro", device.Name())
	a.Equal("UNASPro", device.Model())
	a.NotEmpty(device.SourceName)
	a.NotNil(device.Storage)
	a.NotNil(device.NetworkIO)
	a.Len(device.Drives, 1)
}

func TestUNASGetDevicePartialFailure(t *testing.T) {
	t.Parallel()

	c, requested := newUNASTestClient(t, APIUNASDrivesPath)

	logged := []string{}
	c.ErrorLog = func(msg string, v ...any) {
		logged = append(logged, fmt.Sprintf(msg, v...))
	}

	device, err := c.GetUNASDevice()
	// A console that will not serve one endpoint still reports the rest, and a partial
	// failure is deliberately not an error -- an `if err != nil { return }` at the call
	// site would otherwise discard every metric this console can still supply.
	require.NoError(t, err)
	require.NotNil(t, device)

	a := assert.New(t)
	a.Len(logged, 1, "the endpoint that failed is reported through ErrorLog instead")
	a.Contains(logged[0], "drives")
	a.Len(*requested, unasEndpointCount, "every endpoint is attempted even after one fails")
	a.NotNil(device.DeviceInfo)
	a.NotNil(device.Storage)
	a.NotNil(device.NetworkIO)
	a.Empty(device.Drives)
}

func TestUNASGetDeviceTotalFailure(t *testing.T) {
	t.Parallel()

	c, _ := newUNASTestClient(t,
		APIUNASDeviceInfoPath, APIUNASStoragePath, APIUNASDrivesPath, APIUNASNetworkIOPath)

	device, err := c.GetUNASDevice()
	// Total failure is now the only case that returns an error, and it returns no device.
	require.Error(t, err)
	require.Nil(t, device)
	// The joined error names every endpoint, so an operator sees the whole picture at once.
	require.Contains(t, err.Error(), "device-info")
	require.Contains(t, err.Error(), "storage")
	require.Contains(t, err.Error(), "drives")
	require.Contains(t, err.Error(), "network-io")
}

func TestUNASNilReceivers(t *testing.T) {
	t.Parallel()

	var c *UNASClient

	_, err := c.GetUNASDevice()
	require.ErrorIs(t, err, ErrNilUnifi)

	_, err = (&UNASClient{}).GetUNASDeviceInfo()
	require.ErrorIs(t, err, ErrNilUnifi)

	var d *UNASDevice

	a := assert.New(t)
	a.Empty(d.Name())
	a.Empty(d.Model())
	a.Empty((&UNASDevice{}).Name())
	a.Empty((&UNASDevice{}).Model())
}

func TestNewUNASClientRejectsBadConfig(t *testing.T) {
	t.Parallel()

	_, err := NewUNASClient(nil)
	require.ErrorIs(t, err, ErrNilConfig)

	// Key auth routes Login() through /status, which a storage-only console does not serve.
	_, err = NewUNASClient(&Config{URL: "https://127.0.0.1", APIKey: "abc123"})
	require.ErrorIs(t, err, ErrAPIKeyUnsupported)
}

func TestNewUNASClientLogsIn(t *testing.T) {
	t.Parallel()

	var gotPath, gotMethod, gotBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method

		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)

		w.Header().Set("x-csrf-token", "csrf-abc")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := NewUNASClient(&Config{
		URL: srv.URL, User: "unpoller", Pass: "secret", DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.NoError(t, err)
	require.NotNil(t, c)

	a := assert.New(t)
	// New-style login is assumed rather than probed with a GET of /, which is why no
	// request other than the login lands on the server.
	a.True(c.new)
	a.Equal(APILoginPathNew, gotPath)
	a.Equal(http.MethodPost, gotMethod)
	a.Contains(gotBody, `"username":"unpoller"`)
	a.Equal("csrf-abc", c.csrf)
}

// TestNewUNASClientPinsSSLCert covers the one place NewUNASClient duplicates logic from
// NewUnifi: the fingerprint loop. It writes into a slice newUnifi only allocates when
// SSLCert is non-empty, so an empty-cert regression there is an index panic, not a test
// failure elsewhere. Logging in over TLS also proves the pin actually matches.
func TestNewUNASClientPinsSSLCert(t *testing.T) {
	t.Parallel()

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("x-csrf-token", "csrf-abc")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})
	require.NotNil(t, cert)

	c, err := NewUNASClient(&Config{
		URL: srv.URL, User: "unpoller", Pass: "secret", SSLCert: [][]byte{cert},
		DebugLog: discardLogs, ErrorLog: discardLogs,
	})
	require.NoError(t, err)
	require.NotNil(t, c)

	a := assert.New(t)
	a.Len(c.fingerprints, 1)
	a.NotNil(c.Transport)
	a.Equal("csrf-abc", c.csrf)
}

func TestNewUNASClientLoginFailure(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c, err := NewUNASClient(&Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs})
	require.ErrorIs(t, err, ErrAuthenticationFailed)
	// The client is returned alongside the error so callers can retry a login on it.
	require.NotNil(t, c)
}

func TestUNASStructs(t *testing.T) {
	t.Parallel()

	var info UNASDeviceInfo

	require.NoError(t, gofakeit.Struct(&info))
	require.NotEmpty(t, info.Name)

	var storage UNASStorage

	require.NoError(t, gofakeit.Struct(&storage))

	var drive UNASDrive

	require.NoError(t, gofakeit.Struct(&drive))
	require.NotEmpty(t, drive.ID)

	var netIO UNASNetworkIO

	require.NoError(t, gofakeit.Struct(&netIO))
}
