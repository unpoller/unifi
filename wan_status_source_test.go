package unifi // nolint: testpackage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const wanStatusSourceBody = `{"wan_interfaces":[
	{"name":"Internet 1","state":"ACTIVE","wan_networkgroup":"WAN"},
	{"name":"Cellular","state":"BACKUP","wan_networkgroup":"WAN3"}
]}`

func wanStatusSourceClient(t *testing.T, body string, unifiOS bool) *Unifi {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	return &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    unifiOS,
	}
}

// TestWANStatusCarriesSource is what this change exists for. Every UniFi OS
// console calls its only site "default", so SiteName alone cannot tell two
// controllers apart: a poller watching several of them emits WAN status that
// is indistinguishable, and any downstream attribution has to be guessed.
func TestWANStatusCarriesSource(t *testing.T) {
	t.Parallel()

	c := wanStatusSourceClient(t, wanStatusSourceBody, true)

	status, err := c.GetWANStatus(&Site{Name: "default", SiteName: "Default (default)"})
	require.NoError(t, err)
	require.Len(t, status.WANInterfaces, 2)

	a := assert.New(t)
	a.Equal(c.URL, status.SourceName, "SourceName identifies the controller this came from")
	a.Equal("Default (default)", status.SiteName)
}

// TestWANStatusCarriesSourceOnLegacy pins the same guarantee on the legacy
// path, which classic controllers still take.
func TestWANStatusCarriesSourceOnLegacy(t *testing.T) {
	t.Parallel()

	c := wanStatusSourceClient(t, `{"data":[{"wan_interfaces":[
		{"name":"Internet 1","state":"ACTIVE","wan_networkgroup":"WAN"}
	]}]}`, false)

	status, err := c.GetWANStatus(&Site{Name: "default", SiteName: "Default (default)"})
	require.NoError(t, err)
	require.Len(t, status.WANInterfaces, 1)

	assert.Equal(t, c.URL, status.SourceName)
}

// TestWANStatusCarriesSourceWhenEmpty covers the no-gateway case: the zero
// value returned there must still say where it came from, otherwise a site
// without a gateway is the one entry a caller cannot attribute.
func TestWANStatusCarriesSourceWhenEmpty(t *testing.T) {
	t.Parallel()

	c := wanStatusSourceClient(t, `{}`, true)

	status, err := c.GetWANStatus(&Site{Name: "default", SiteName: "Default (default)"})
	require.NoError(t, err)
	assert.Empty(t, status.WANInterfaces)
	assert.Equal(t, c.URL, status.SourceName)
}
