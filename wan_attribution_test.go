package unifi // nolint: testpackage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The enriched-configuration endpoint answers with a bare array and carries no
// site or controller identity of its own — which is exactly why callers cannot
// tell two controllers apart without help.
const wanEnrichedBody = `[
	{"configuration":{"_id":"a1","name":"Internet 1","wan_networkgroup":"WAN"}},
	{"configuration":{"_id":"a2","name":"Internet 2","wan_networkgroup":"WAN2"}}
]`

// TestWANEnrichedConfigurationIsAttributed pins the fix: every entry must come
// back tagged with the site it was fetched for. Without it, a poller watching
// several controllers emits WAN metrics that are indistinguishable, and any
// downstream label has to be guessed.
func TestWANEnrichedConfigurationIsAttributed(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(wanEnrichedBody))
	}))
	t.Cleanup(srv.Close)

	c := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
	}

	site := &Site{Name: "default", SiteName: "Default (default)", SourceName: srv.URL}

	wans, err := c.GetWANEnrichedConfiguration([]*Site{site})
	require.NoError(t, err)
	require.Len(t, wans, 2)

	a := assert.New(t)
	for _, w := range wans {
		a.Equal("Default (default)", w.SiteName)
		a.Equal(srv.URL, w.SourceName)
	}
	a.Equal("Internet 1", wans[0].Configuration.Name)
	a.Equal("WAN2", wans[1].Configuration.WANNetworkgroup)
}

// TestWANEnrichedConfigurationPerSite covers the multi-site case, the one the
// fix exists for: each entry keeps the identity of the site it came from
// instead of inheriting whichever was fetched last.
func TestWANEnrichedConfigurationPerSite(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/proxy/network/v2/api/site/second/wan/enriched-configuration" {
			_, _ = w.Write([]byte(`[{"configuration":{"_id":"b1","name":"Fiber","wan_networkgroup":"WAN"}}]`))

			return
		}

		_, _ = w.Write([]byte(`[{"configuration":{"_id":"a1","name":"Internet 1","wan_networkgroup":"WAN"}}]`))
	}))
	t.Cleanup(srv.Close)

	c := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
		new:    true,
	}

	wans, err := c.GetWANEnrichedConfiguration([]*Site{
		{Name: "default", SiteName: "First (default)", SourceName: "https://one"},
		{Name: "second", SiteName: "Second (second)", SourceName: "https://two"},
	})
	require.NoError(t, err)
	require.Len(t, wans, 2)

	a := assert.New(t)
	a.Equal("First (default)", wans[0].SiteName)
	a.Equal("https://one", wans[0].SourceName)
	a.Equal("Second (second)", wans[1].SiteName)
	a.Equal("https://two", wans[1].SourceName)
}
