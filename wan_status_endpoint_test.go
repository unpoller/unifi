package unifi // nolint: testpackage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Payload observed on a UDR running UniFi OS: the v2 load-balancing endpoint
// answers with the interfaces at the top level, no "data" envelope.
const wanStatusV2Body = `{"wan_interfaces":[
	{"name":"Internet 1","state":"ACTIVE","wan_networkgroup":"WAN"},
	{"name":"Internet 2","state":"BACKUP","wan_networkgroup":"WAN2"},
	{"name":"Cellular","state":"BACKUP","wan_networkgroup":"WAN3"}
]}`

// The legacy endpoint wraps the same object in a "data" array.
const wanStatusLegacyBody = `{"data":[{"wan_interfaces":[
	{"name":"Internet 1","state":"ACTIVE","wan_networkgroup":"WAN"}
]}]}`

func wanStatusServer(t *testing.T, body string) (*httptest.Server, *[]string) {
	t.Helper()

	requested := &[]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requested = append(*requested, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	return srv, requested
}

func wanStatusClient(srv *httptest.Server, unifiOS bool) *Unifi {
	return &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    unifiOS,
	}
}

// TestGetWANStatusUnifiOS is the regression this change exists for: a UniFi OS
// console answers the legacy /stat/status endpoint with no interfaces at all —
// no error, just an empty result — so every caller saw a site with no WAN. The
// interfaces must be read from the v2 load-balancing endpoint instead.
func TestGetWANStatusUnifiOS(t *testing.T) {
	t.Parallel()

	srv, requested := wanStatusServer(t, wanStatusV2Body)
	c := wanStatusClient(srv, true)

	status, err := c.GetWANStatus(&Site{Name: "default", SiteName: "Default (default)"})
	require.NoError(t, err)
	require.Len(t, status.WANInterfaces, 3)

	a := assert.New(t)
	a.Equal([]string{"/proxy/network/v2/api/site/default/wan/load-balancing/status"}, *requested)
	a.Equal("Internet 1", status.WANInterfaces[0].Name)
	a.Equal("ACTIVE", status.WANInterfaces[0].State)
	a.Equal("WAN", status.WANInterfaces[0].WANNetworkgroup)
	a.Equal("BACKUP", status.WANInterfaces[2].State)
	// SiteName carries json:"-", so it survives the unmarshal and must be set.
	a.Equal("Default (default)", status.SiteName)
}

// TestGetWANStatusClassic pins the old behaviour: a classic controller has no
// v2 endpoint, so the legacy path and its "data" envelope must keep working.
func TestGetWANStatusClassic(t *testing.T) {
	t.Parallel()

	srv, requested := wanStatusServer(t, wanStatusLegacyBody)
	c := wanStatusClient(srv, false)

	status, err := c.GetWANStatus(&Site{Name: "default", SiteName: "Default (default)"})
	require.NoError(t, err)
	require.Len(t, status.WANInterfaces, 1)

	a := assert.New(t)
	a.Equal([]string{"/api/s/default/stat/status"}, *requested)
	a.Equal("ACTIVE", status.WANInterfaces[0].State)
	a.Equal("Default (default)", status.SiteName)
}

// TestGetWANStatusUnifiOSNoGateway covers a site with no gateway: the endpoint
// answers with an empty object. That is not an error, and callers detect it by
// checking len(status.WANInterfaces) == 0, as the doc comment promises.
func TestGetWANStatusUnifiOSNoGateway(t *testing.T) {
	t.Parallel()

	srv, _ := wanStatusServer(t, `{}`)
	c := wanStatusClient(srv, true)

	status, err := c.GetWANStatus(&Site{Name: "default", SiteName: "Default (default)"})
	require.NoError(t, err)
	assert.Empty(t, status.WANInterfaces)
	assert.Equal(t, "Default (default)", status.SiteName)
}
