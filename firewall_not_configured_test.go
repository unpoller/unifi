package unifi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

const zoneBasedFirewallNotConfiguredBody = `{"statusCode":400,"statusName":"BAD_REQUEST","code":"api.firewall.zone-based-firewall-not-configured","message":"Zone Based Firewall is not configured","requestPath":"/integration/v1/sites/site-id/firewall/zones"}`

func TestGetFirewallZonesNotConfigured(t *testing.T) {
	t.Parallel()

	u, logged := firewallClient(t, http.StatusBadRequest, zoneBasedFirewallNotConfiguredBody)

	zones, err := u.GetFirewallZones(&unifi.IntegrationSite{ID: "site-id", Name: "default"})
	require.NoError(t, err)
	assert.Empty(t, zones)
	assert.Empty(t, *logged, "an unconfigured zone-based firewall is not an error")
}

func TestGetFirewallZonesOtherBadRequest(t *testing.T) {
	t.Parallel()

	body := `{"statusCode":400,"statusName":"BAD_REQUEST","code":"api.firewall.something-else","message":"nope"}`
	u, logged := firewallClient(t, http.StatusBadRequest, body)

	_, err := u.GetFirewallZones(&unifi.IntegrationSite{ID: "site-id", Name: "default"})
	require.Error(t, err)
	assert.ErrorIs(t, err, unifi.ErrInvalidStatusCode)
	assert.NotEmpty(t, *logged)
}

func TestGetFirewallPoliciesNotConfigured(t *testing.T) {
	t.Parallel()

	u, logged := firewallClient(t, http.StatusBadRequest, zoneBasedFirewallNotConfiguredBody)

	policies, err := u.GetFirewallPolicies([]*unifi.Site{{Name: "default", SiteName: "Default"}})
	require.NoError(t, err)
	assert.Empty(t, policies)
	assert.Empty(t, *logged)
}

func TestGetFirewallPoliciesOtherBadRequest(t *testing.T) {
	t.Parallel()

	u, logged := firewallClient(t, http.StatusBadRequest, `{"meta":{"rc":"error","msg":"api.err.Invalid"}}`)

	_, err := u.GetFirewallPolicies([]*unifi.Site{{Name: "default", SiteName: "Default"}})
	require.Error(t, err)
	assert.ErrorIs(t, err, unifi.ErrInvalidStatusCode)
	assert.Empty(t, *logged, "policy failures are returned to the caller; the library does not log them")
}

func firewallClient(t *testing.T, status int, body string) (*unifi.Unifi, *[]string) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	logged := []string{}
	u := &unifi.Unifi{
		Client: srv.Client(),
		Config: &unifi.Config{
			URL:      srv.URL,
			APIKey:   "test-key",
			DebugLog: func(string, ...interface{}) {},
			ErrorLog: func(msg string, args ...interface{}) {
				logged = append(logged, fmt.Sprintf(msg, args...))
			},
		},
	}

	return u, &logged
}
