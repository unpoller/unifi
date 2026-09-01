package unifi // nolint: testpackage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// integrationSourceClient serves the same one-element payload to every
// Integration/v1 endpoint. The shape does not matter here: what is under test
// is the attribution the library stamps on its way out, not the decoding.
func integrationSourceClient(t *testing.T) *Unifi {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"a1","name":"one"}],"totalCount":1,"count":1}`))
	}))
	t.Cleanup(srv.Close)

	return &Unifi{
		Client: &http.Client{},
		Config: &Config{
			URL:      srv.URL,
			APIKey:   "test-key",
			DebugLog: discardLogs,
			ErrorLog: discardLogs,
		},
		new: true,
	}
}

// TestIntegrationEntitiesCarrySource is what this change exists for. These
// families carried SiteName alone, and SiteName is "default" on every UniFi OS
// console — so a poller watching several controllers produced entities that
// could not be told apart. Each getter must stamp the controller it read from.
func TestIntegrationEntitiesCarrySource(t *testing.T) {
	t.Parallel()

	c := integrationSourceClient(t)
	site := &IntegrationSite{ID: "site-id", Name: "default"}

	t.Run("VPNServer", func(t *testing.T) {
		got, err := c.GetVPNServers(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("SiteToSiteTunnel", func(t *testing.T) {
		got, err := c.GetSiteToSiteTunnels(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("FirewallZone", func(t *testing.T) {
		got, err := c.GetFirewallZones(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("ACLRule", func(t *testing.T) {
		got, err := c.GetACLRules(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("WifiBroadcast", func(t *testing.T) {
		got, err := c.GetWifiBroadcasts(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("DNSPolicy", func(t *testing.T) {
		got, err := c.GetDNSPolicies(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("RADIUSProfile", func(t *testing.T) {
		got, err := c.GetRADIUSProfiles(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("TrafficMatchingList", func(t *testing.T) {
		got, err := c.GetTrafficMatchingLists(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})

	t.Run("HotspotVoucher", func(t *testing.T) {
		got, err := c.GetHotspotVouchers(site)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, c.URL, got[0].SourceName)
	})
}
