package unifi // nolint: testpackage

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizeGuest(t *testing.T) {
	t.Parallel()

	type stamgrBody struct {
		Cmd     string `json:"cmd"`
		MAC     string `json:"mac"`
		Minutes int    `json:"minutes,omitempty"`
	}

	var got stamgrBody

	var gotPath, gotMethod string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err := json.Unmarshal(body, &got); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[],"meta":{"rc":"ok"}}`))
	}))
	defer srv.Close()

	u := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
	}
	site := &Site{Name: "default"}

	err := u.AuthorizeGuest(site, "aa:bb:cc:dd:ee:ff", 60)
	require.NoError(t, err)

	a := assert.New(t)
	a.Equal("/api/s/default/cmd/stamgr", gotPath)
	a.Equal(http.MethodPost, gotMethod)
	a.Equal("authorize-guest", got.Cmd)
	a.Equal("aa:bb:cc:dd:ee:ff", got.MAC)
	a.Equal(60, got.Minutes)
}

func TestAuthorizeGuestZeroMinutesOmitted(t *testing.T) {
	t.Parallel()

	var raw map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err := json.Unmarshal(body, &raw); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[],"meta":{"rc":"ok"}}`))
	}))
	defer srv.Close()

	u := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
	}

	err := u.AuthorizeGuest(&Site{Name: "default"}, "aa:bb:cc:dd:ee:ff", 0)
	require.NoError(t, err)

	a := assert.New(t)
	_, hasMinutes := raw["minutes"]
	a.False(hasMinutes, "minutes should be omitted when 0")
}

func TestAuthorizeGuestValidation(t *testing.T) {
	t.Parallel()

	u := &Unifi{Client: &http.Client{}, Config: &Config{URL: "http://invalid", DebugLog: discardLogs, ErrorLog: discardLogs}}

	a := assert.New(t)
	a.ErrorIs(u.AuthorizeGuest(nil, "aa:bb:cc:dd:ee:ff", 0), ErrNoSiteProvided)
	a.ErrorIs(u.AuthorizeGuest(&Site{Name: ""}, "aa:bb:cc:dd:ee:ff", 0), ErrNoSiteProvided)
	a.ErrorIs(u.AuthorizeGuest(&Site{Name: "default"}, "", 0), ErrEmptyMAC)
	a.ErrorIs(u.AuthorizeGuest(&Site{Name: "default"}, "aa:bb:cc:dd:ee:ff", -1), ErrInvalidMinutes)
}

func TestAuthorizeGuestControllerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"data":[],"meta":{"rc":"error","msg":"api.err.LoginRequired"}}`))
	}))
	defer srv.Close()

	u := &Unifi{
		Client: &http.Client{},
		Config: &Config{URL: srv.URL, DebugLog: discardLogs, ErrorLog: discardLogs},
	}

	err := u.AuthorizeGuest(&Site{Name: "default"}, "aa:bb:cc:dd:ee:ff", 0)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, ErrInvalidStatusCode), "error should unwrap to ErrInvalidStatusCode")
	a.Contains(err.Error(), "aa:bb:cc:dd:ee:ff", "error should include MAC for context")
}
