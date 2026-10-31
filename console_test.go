package unifi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// Reduced response from unpoller/unpoller#1130; other cases below are synthetic.
const consoleUpdateBody = `{"devices":{"unifiOS":[{
	"name":"UCG Fiber","firmware":"5.1.33",
	"firmwareStatus":"updateAvailable","updateAvailable":"5.1.34+4d109cd"
}]},"apps":{"controllers":[{"name":"network","updateAvailable":null}]}}`

func consoleClient(t *testing.T, legacy bool, handler http.HandlerFunc) *unifi.Unifi {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc(unifi.APIConsoleSystemPath, handler)
	mux.HandleFunc(unifi.APILoginPathNew, func(w http.ResponseWriter, _ *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: "test-session", Path: "/"})
		w.Header().Set("X-CSRF-Token", "test-csrf")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if legacy && r.URL.Path == "/" {
			w.Header().Set("Location", "/manage")
			w.WriteHeader(http.StatusFound)

			return
		}

		_, err := w.Write([]byte(`{"meta":{"server_version":"10.6.109"}}`))
		assert.NoError(t, err)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client, err := unifi.NewUnifi(&unifi.Config{
		URL: srv.URL, User: "monitor", Pass: "test-password", Timeout: time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(client.CloseIdleConnections)

	return client
}

func TestGetConsoleSystem(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32

	client := consoleClient(t, false, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test-csrf", r.Header.Get("X-CSRF-Token"))

		cookie, err := r.Cookie("TOKEN")
		if assert.NoError(t, err) {
			assert.Equal(t, "test-session", cookie.Value)
		}

		_, err = w.Write([]byte(consoleUpdateBody))
		assert.NoError(t, err)
	})

	got, err := client.GetConsoleSystem()
	require.NoError(t, err)
	require.Len(t, got.Devices.UniFiOS, 1)
	assert.EqualValues(t, 1, requests.Load())
	assert.Equal(t, client.URL, got.SourceName)
	assert.Equal(t, "UCG Fiber", got.Devices.UniFiOS[0].Name)
	assert.Equal(t, "5.1.33", got.Devices.UniFiOS[0].Firmware)
	assert.Equal(t, "5.1.34+4d109cd", got.Devices.UniFiOS[0].UpdateAvailable)
	available, known := got.Devices.UniFiOS[0].UpdateStatus()
	assert.True(t, available)
	assert.True(t, known)

	// Exempt only the console endpoint; existing Network and Protect routing stays intact.
	for path, expected := range map[string]string{
		unifi.APIConsoleSystemPath:    "/api/system",
		"/api/s/default/stat/sysinfo": "/proxy/network/api/s/default/stat/sysinfo",
		"/proxy/protect/api/events":   "/proxy/protect/api/events",
		unifi.APILoginPath:            unifi.APILoginPathNew,
	} {
		req, err := client.UniReq(path, "")
		require.NoError(t, err)
		assert.Equal(t, expected, req.URL.Path)
	}
}

func TestGetConsoleSystemUnavailable(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   error
	}{
		{"not found", 404, `{}`, unifi.ErrEndpointNotFound},
		{"forbidden", 403, `{}`, unifi.ErrInvalidStatusCode},
		{"unauthorized", 401, `{}`, unifi.ErrInvalidStatusCode},
		{"rate limited", 429, `{}`, unifi.ErrTooManyRequests},
		{"missing devices", 200, `{}`, unifi.ErrNoConsoleSystemData},
		{"empty array", 200, `{"devices":{"unifiOS":[]}}`, unifi.ErrNoConsoleSystemData},
		{"null entry", 200, `{"devices":{"unifiOS":[null]}}`, unifi.ErrNoConsoleSystemData},
		{"null response", 200, `null`, unifi.ErrNoConsoleSystemData},
		{"HTML", 200, `<html></html>`, nil},
		{"wrong version type", 200, `{"devices":{"unifiOS":[{"updateAvailable":true}]}}`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := consoleClient(t, false, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, err := w.Write([]byte(tc.body))
				assert.NoError(t, err)
			})
			got, err := client.GetConsoleSystem()
			require.Error(t, err)
			assert.Nil(t, got)

			if tc.want != nil {
				assert.ErrorIs(t, err, tc.want)
			}
		})
	}
}

func TestGetConsoleSystemLegacyAndNil(t *testing.T) {
	t.Parallel()

	client := consoleClient(t, true, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("legacy controllers must not request the console endpoint")
	})
	got, err := client.GetConsoleSystem()
	require.ErrorIs(t, err, unifi.ErrConsoleSystemUnsupported)
	assert.Nil(t, got)

	got, err = (*unifi.Unifi)(nil).GetConsoleSystem()
	require.ErrorIs(t, err, unifi.ErrNilUnifi)
	assert.Nil(t, got)
}

func TestConsoleFirmwareUpdateStatus(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		body      string
		available bool
		known     bool
	}{
		{"offered version", `{"updateAvailable":"5.1.34"}`, true, true},
		{"available status", `{"firmwareStatus":"updateAvailable"}`, true, true},
		{"up to date", `{"firmwareStatus":"upToDate","updateAvailable":null}`, false, true},
		{"up to date empty version", `{"firmwareStatus":"upToDate","updateAvailable":""}`, false, true},
		{"missing", `{}`, false, false},
		{"null version only", `{"updateAvailable":null}`, false, false},
		{"empty version only", `{"updateAvailable":""}`, false, false},
		{"unknown status", `{"firmwareStatus":"checking"}`, false, false},
		{"offer takes precedence", `{"firmwareStatus":"upToDate","updateAvailable":"5.1.34"}`, true, true},
		{"nil", `null`, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var firmware *unifi.ConsoleFirmware

			require.NoError(t, json.Unmarshal([]byte(tc.body), &firmware))
			available, known := firmware.UpdateStatus()
			assert.Equal(t, tc.available, available)
			assert.Equal(t, tc.known, known)
		})
	}
}
