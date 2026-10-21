package unifi // nolint: testpackage

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skipTagsClient serves an empty device list and counts how often the
// device-tags endpoint is asked for, which is the only thing under test here.
func skipTagsClient(t *testing.T, skip bool) (*Unifi, *atomic.Int32) {
	t.Helper()

	var tagCalls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.HasSuffix(r.URL.Path, "/device-tags") {
			tagCalls.Add(1)

			_, _ = w.Write([]byte(`[]`))

			return
		}

		_, _ = w.Write([]byte(`{"data":[],"meta":{"rc":"ok"}}`))
	}))
	t.Cleanup(srv.Close)

	return &Unifi{
		Client: &http.Client{},
		Config: &Config{
			URL:            srv.URL,
			APIKey:         "test-key",
			DebugLog:       discardLogs,
			ErrorLog:       discardLogs,
			SkipDeviceTags: skip,
		},
		new: true,
	}, &tagCalls
}

func TestGetDevicesSkipDeviceTags(t *testing.T) {
	t.Parallel()

	sites := []*Site{{Name: "default", SiteName: "Default"}}

	t.Run("default fetches tags once per site", func(t *testing.T) {
		t.Parallel()

		u, tagCalls := skipTagsClient(t, false)

		devices, err := u.GetDevices(sites)
		require.NoError(t, err)
		require.NotNil(t, devices)
		assert.Equal(t, int32(1), tagCalls.Load())
	})

	t.Run("SkipDeviceTags never calls the endpoint", func(t *testing.T) {
		t.Parallel()

		u, tagCalls := skipTagsClient(t, true)

		devices, err := u.GetDevices(sites)
		require.NoError(t, err)
		require.NotNil(t, devices)
		assert.Equal(t, int32(0), tagCalls.Load())
	})
}
