package unifi // nolint: testpackage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnifiNilConfig(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	uni, err := NewUnifi(nil)
	a.Nil(uni)
	a.Error(err)
	a.Contains(err.Error(), "config is nil")
}

func TestNewUnifi(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u := "http://127.0.0.1:64431"
	c := &Config{
		User:      "user1",
		Pass:      "pass2",
		URL:       u,
		VerifySSL: false,
		DebugLog:  discardLogs,
	}
	authReq, err := NewUnifi(c)
	a.NotNil(err)
	a.EqualValues(u, authReq.URL)
	a.Contains(err.Error(), "connection refused", "an invalid destination should produce a connection error.")
}

func TestNewUnifiAPIKey(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u := "http://127.0.0.1:64431"
	c := &Config{
		APIKey:    "fakekey",
		URL:       u,
		VerifySSL: false,
		DebugLog:  discardLogs,
	}
	authReq, err := NewUnifi(c)
	a.NotNil(err)
	a.EqualValues(u, authReq.URL)
	a.Contains(err.Error(), "connection refused", "an invalid destination should produce a connection error.")
}

func TestUniReq(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := "/test/path"
	u := "http://some.url:8443"

	// Test empty parameters.
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{URL: u, DebugLog: discardLogs}}
	r, err := authReq.UniReq(p, "")

	a.Nil(err, "newrequest must not produce an error")
	a.EqualValues(p, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(u, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("GET", r.Method, "without parameters the method must be GET")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Test with parameters
	k := "key1=value9&key2=value7"
	authReq = &Unifi{Client: &http.Client{}, Config: &Config{URL: "http://some.url:8443", DebugLog: discardLogs}}
	r, err = authReq.UniReq(p, k)
	a.Nil(err, "newrequest must not produce an error")

	a.EqualValues(p, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(u, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("POST", r.Method, "with parameters the method must be POST")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Check the parameters.
	d, err := io.ReadAll(r.Body)
	a.Nil(err, "problem reading request body, POST parameters may be malformed")
	a.EqualValues(k, string(d), "POST parameters improperly encoded")
}

func TestUniReqPut(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := "/test/path"
	u := "http://some.url:8443"

	// Test empty parameters.
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{URL: u, DebugLog: discardLogs}}
	_, err := authReq.UniReqPut(p, "")
	a.NotNil(err, "empty params must produce an error")

	// Test with parameters
	k := "key1=value9&key2=value7"
	authReq = &Unifi{Client: &http.Client{}, Config: &Config{URL: "http://some.url:8443", DebugLog: discardLogs}}
	r, err := authReq.UniReqPut(p, k)
	a.Nil(err, "newrequest must not produce an error")

	a.EqualValues(p, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(u, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("PUT", r.Method, "with parameters the method must be POST")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Check the parameters.
	d, err := io.ReadAll(r.Body)
	a.Nil(err, "problem reading request body, PUT parameters may be malformed")
	a.EqualValues(k, string(d), "PUT parameters improperly encoded")
}

func TestUnifiIntegrationAPIKeyInjected(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") == "fakekey" {
			w.WriteHeader(http.StatusOK)

			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{APIKey: "fakekey", URL: srv.URL, DebugLog: discardLogs}}
	authResp, err := authReq.UniReqPost("/test", "")
	a.Nil(err, "newrequest must not produce an error")
	a.EqualValues("POST", authResp.Method, "with parameters the method must be POST")
}

func TestUnifiIntegrationUserPassInjected(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.URL.Path, "/api/login") {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("error reading body:%v\n", err)

			return
		}

		type userPass struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var up userPass

		err = json.Unmarshal(data, &up)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("error decoding body: %s: %s\n", string(data), err)

			return
		}

		if strings.EqualFold(up.Username, "fakeuser") && strings.EqualFold(up.Password, "fakepass") {
			w.WriteHeader(http.StatusOK)
		}

		w.WriteHeader(http.StatusUnauthorized)
	}))
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{User: "fakeuser", Pass: "fakepass", URL: srv.URL, DebugLog: discardLogs}}
	err := authReq.Login()
	a.Nil(err, "user/pass login must not produce an error")
}

func TestParseRetryAfter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	a.Equal(60*time.Second, parseRetryAfter(""))
	a.Equal(120*time.Second, parseRetryAfter("120"))
	a.Equal(60*time.Second, parseRetryAfter("0")) // 0 not > 0, so default
	a.Equal(5*time.Minute, parseRetryAfter("9999"))
	a.Equal(60*time.Second, parseRetryAfter("invalid"))
}

func TestRateLimitError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	e := &RateLimitError{RetryAfter: 90 * time.Second}
	a.True(errors.Is(e, ErrTooManyRequests))
	a.Contains(e.Error(), "429")
	a.Contains(e.Error(), "retry after")
}

/* NOT DONE: OPEN web server, check parameters posted, more. These tests are incomplete.
a.EqualValues(`{"username": "user1","password": "pass2"}`, string(post_params),
	"user/pass json parameters improperly encoded")
*/

// TestDo_GzipBodyWithoutContentEncoding mirrors unpoller/unpoller#997 on the
// local controller path: a proxy returns a gzipped JSON body but does not set
// Content-Encoding, so Go's transport hands raw compressed bytes to the
// client. The library must detect and decompress the body itself.
func TestDo_GzipBodyWithoutContentEncoding(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"data":[{"name":"default"}],"meta":{"rc":"ok"}}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Intentionally do NOT set Content-Encoding: gzip.
		_, _ = w.Write(gzipBytes(t, payload))
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	got, err := u.GetJSON("/api/s/default/self/sites")
	require.NoError(t, err)
	assert.Equal(t, payload, got)
}

// truncatedGzip returns a valid gzip header followed by a cut-off deflate
// stream that fails to decode. Used to simulate a corrupt-but-magic-prefixed
// body on error responses.
func truncatedGzip() []byte {
	return []byte{0x1f, 0x8b, 0x08, 0x00, 0x00}
}

// TestDo_GzippedRateLimit_PreservesStatus locks in that decompression failures
// on error responses do not mask the HTTP status. A 429 with a malformed gzip
// body must still surface as *RateLimitError, and the raw (undecoded) body
// must be returned alongside the error so callers can still log it.
func TestDo_GzippedRateLimit_PreservesStatus(t *testing.T) {
	t.Parallel()

	raw := truncatedGzip()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write(raw)
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	got, err := u.GetJSON("/api/s/default/self/sites")
	require.Error(t, err)

	var rateErr *RateLimitError
	assert.ErrorAs(t, err, &rateErr, "gzip decode failure must not mask 429 status")
	assert.Equal(t, raw, got, "raw body must be returned when gzip decode fails on error responses")
}

// TestDo_GzippedNotFound_PreservesStatus ensures a 404 with a malformed gzip
// body still surfaces ErrEndpointNotFound rather than the generic decoding
// error, since callers (e.g. endpoint probes) rely on the sentinel.
func TestDo_GzippedNotFound_PreservesStatus(t *testing.T) {
	t.Parallel()

	raw := truncatedGzip()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(raw)
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	got, err := u.GetJSON("/api/s/default/self/sites")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEndpointNotFound, "decompression failure must not hide 404")
	assert.Equal(t, raw, got, "raw body must be returned when gzip decode fails on error responses")
}

// TestDo_GzippedServerError_PreservesStatus covers the generic non-2xx,
// non-404, non-429 branch: a 500 with a malformed gzip body must still surface
// ErrInvalidStatusCode so callers checking the sentinel still work.
func TestDo_GzippedServerError_PreservesStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(truncatedGzip())
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	_, err := u.GetJSON("/api/s/default/self/sites")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatusCode, "decompression failure must not hide 500")
	assert.Contains(t, err.Error(), "500")
}

// TestDo_GzippedSuccess_DecompressionFailureIsFatal ensures a 2xx response
// with a malformed gzip body is reported as an error (and not silently passed
// through as raw bytes to the JSON parser).
func TestDo_GzippedSuccess_DecompressionFailureIsFatal(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(truncatedGzip())
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	_, err := u.GetJSON("/api/s/default/self/sites")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decoding response")
}

// TestDo_GzippedServerError_ValidGzip_ReturnsDecompressedBody covers the
// remaining branch where decompression *succeeds* on a non-2xx response: the
// status-aware sentinel must still fire, and the caller must receive the
// decompressed (not the raw gzipped) body so error logging is readable.
func TestDo_GzippedServerError_ValidGzip_ReturnsDecompressedBody(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"meta":{"rc":"error","msg":"api.err.ServerError"}}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(gzipBytes(t, payload))
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	got, err := u.GetJSON("/api/s/default/self/sites")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatusCode)
	assert.Equal(t, payload, got, "decompressed body must be returned when gzip decode succeeds, even on error status")
}

// TestDo_PlainJSONBody is the non-gzipped baseline so we know the decompress
// hook does not interfere with normal responses.
func TestDo_PlainJSONBody(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"data":[{"name":"default"}],"meta":{"rc":"ok"}}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	u := &Unifi{Client: srv.Client(), Config: &Config{URL: srv.URL, DebugLog: discardLogs}}

	got, err := u.GetJSON("/api/s/default/self/sites")
	require.NoError(t, err)
	assert.Equal(t, payload, got)
}
