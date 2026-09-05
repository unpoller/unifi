package unifi // nolint: testpackage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The controller wraps sysinfo in the usual {"meta":…,"data":[…]} envelope.
const sysinfoBody = `{"meta":{"rc":"ok"},"data":[{
	"timezone":"Etc/UTC",
	"autobackup":false,
	"build":"atag_9.9.9_00000",
	"version":"9.9.9",
	"previous_version":"9.9.8",
	"data_retention_days":90,
	"data_retention_time_in_hours_for_5minutes_scale":24,
	"update_available":false,
	"hostname":"unifi",
	"name":"Console",
	"inform_port":8080,
	"https_port":8443,
	"uptime":123456,
	"has_webrtc_support":true,
	"ubnt_device_type":"UCKP",
	"unsupported_device_count":0,
	"is_cloud_console":false,
	"console_display_version":"9.0.0"
}]}`

const sysinfoEmptyBody = `{"meta":{"rc":"ok"},"data":[]}`

func sysinfoClient(t *testing.T, body string) (*Unifi, *[]string) {
	t.Helper()

	requested := &[]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requested = append(*requested, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	return &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
	}, requested
}

// TestGetSysinfoSite is the regression this change exists for: the response was
// unmarshalled into a bare Sysinfo, so none of the fields inside the "data"
// envelope matched. encoding/json reports no error for that, so every caller
// received a fully zeroed struct and a nil error.
func TestGetSysinfoSite(t *testing.T) {
	t.Parallel()

	c, requested := sysinfoClient(t, sysinfoBody)

	s, err := c.GetSysinfoSite(&Site{Name: "default", SiteName: "Default (default)", SourceName: "controller"})
	require.NoError(t, err)
	require.NotNil(t, s)

	a := assert.New(t)
	a.Equal([]string{"/api/s/default/stat/sysinfo"}, *requested)
	a.Equal("9.9.9", s.Version)
	a.Equal("atag_9.9.9_00000", s.Build)
	a.Equal("9.9.8", s.PreviousVer)
	a.Equal("9.0.0", s.ConsoleVer)
	a.Equal("UCKP", s.DeviceType)
	a.Equal("Console", s.Name)
	a.Equal(90, s.DataRetDays)
	a.Equal(24, s.DataRet5min)
	a.Equal(8443, s.HTTPSPort)
	a.Equal(int64(123456), s.Uptime)
	a.True(s.HasWebRTC)
	// SiteName and SourceName carry json:"-" and must be filled in by the caller.
	a.Equal("Default (default)", s.SiteName)
	a.Equal("controller", s.SourceName)
}

// TestGetSysinfoSiteEmpty pins the other half: an empty data array must surface
// as an error rather than as a zeroed struct that looks like real data.
func TestGetSysinfoSiteEmpty(t *testing.T) {
	t.Parallel()

	c, _ := sysinfoClient(t, sysinfoEmptyBody)

	s, err := c.GetSysinfoSite(&Site{Name: "default"})
	require.ErrorIs(t, err, ErrNoSysinfoData)
	assert.Nil(t, s)
}

// TestGetSysinfo walks the per-site loop.
func TestGetSysinfo(t *testing.T) {
	t.Parallel()

	c, _ := sysinfoClient(t, sysinfoBody)

	all, err := c.GetSysinfo([]*Site{{Name: "default", SiteName: "Default (default)"}})
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "9.9.9", all[0].Version)
}
