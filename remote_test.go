package unifi // nolint: testpackage

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func gzipBytes(t *testing.T, payload []byte) []byte {
	t.Helper()

	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)
	_, err := w.Write(payload)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	return buf.Bytes()
}

func TestMaybeDecompressGzip_NoMagic_PassesThrough(t *testing.T) {
	t.Parallel()

	in := []byte(`{"data":[]}`)

	out, err := maybeDecompressGzip(in)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestMaybeDecompressGzip_ShortBody_PassesThrough(t *testing.T) {
	t.Parallel()

	out, err := maybeDecompressGzip([]byte{0x1f})
	require.NoError(t, err)
	assert.Equal(t, []byte{0x1f}, out)
}

func TestMaybeDecompressGzip_GzipMagic_Decompresses(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"data":[{"id":"abc","name":"default"}]}`)
	compressed := gzipBytes(t, payload)

	out, err := maybeDecompressGzip(compressed)
	require.NoError(t, err)
	assert.Equal(t, payload, out)
}

func TestMaybeDecompressGzip_InvalidGzip_ReturnsError(t *testing.T) {
	t.Parallel()

	// Magic bytes present but body is otherwise invalid gzip.
	bad := []byte{0x1f, 0x8b, 0x00, 0x00}

	_, err := maybeDecompressGzip(bad)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gzip")
}

func TestMaybeDecompressGzip_EmptyBody_PassesThrough(t *testing.T) {
	t.Parallel()

	out, err := maybeDecompressGzip(nil)
	require.NoError(t, err)
	assert.Empty(t, out)

	out, err = maybeDecompressGzip([]byte{})
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestMaybeDecompressGzip_TruncatedAfterHeader_ReturnsError(t *testing.T) {
	t.Parallel()

	full := gzipBytes(t, []byte(`{"data":[]}`))
	// Keep the magic+header but truncate the deflate stream.
	truncated := full[:10]

	_, err := maybeDecompressGzip(truncated)
	require.Error(t, err)
}

func TestMaybeDecompressGzip_CRCMismatch_ReturnsError(t *testing.T) {
	t.Parallel()

	full := gzipBytes(t, []byte(`{"data":[]}`))
	// Corrupt the CRC32 in the trailer (last 8 bytes are CRC32 + ISIZE).
	corrupted := append([]byte(nil), full...)
	corrupted[len(corrupted)-8] ^= 0xff

	_, err := maybeDecompressGzip(corrupted)
	require.Error(t, err)
}

// TestDiscoverSites_GzipBodyWithoutContentEncoding reproduces unpoller/unpoller#997:
// the Site Manager proxy returns a gzipped body but does not set
// Content-Encoding: gzip, so Go's transport hands raw compressed bytes to the
// client. The library must detect and decompress the body itself.
func TestDiscoverSites_GzipBodyWithoutContentEncoding(t *testing.T) {
	t.Parallel()

	resp := SitesResponse{
		Data: []RemoteSite{
			{ID: "site-1", InternalReference: "default", Name: "Default"},
		},
		HTTPStatusCode: 200,
		TraceID:        "trace-xyz",
	}

	payload, err := json.Marshal(resp)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Intentionally do NOT set Content-Encoding: gzip — this mirrors the
		// upstream Site Manager behavior described in the bug report.
		_, _ = w.Write(gzipBytes(t, payload))
	}))
	defer srv.Close()

	client := NewRemoteAPIClient("test-key", nil, nil, nil)
	client.baseURL = srv.URL

	sites, err := client.DiscoverSites("console-1")
	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, "site-1", sites[0].ID)
	assert.Equal(t, "default", sites[0].InternalReference)
	assert.Equal(t, "Default", sites[0].Name)
}

// TestMakeRequestOnce_GzippedRateLimit_PreservesStatus locks in that
// decompression failures on error responses do not mask the HTTP status. A 429
// with a malformed gzip body must still surface as *RateLimitError so
// makeRequest can honor Retry-After. The 64 KiB LimitReader on error bodies
// makes truncated gzip streams a realistic failure mode.
//
// Calls makeRequestOnce directly to skip the retry loop's Retry-After sleep.
func TestMakeRequestOnce_GzippedRateLimit_PreservesStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
		// Magic prefix + truncated stream — would fail gzip decode.
		_, _ = w.Write([]byte{0x1f, 0x8b, 0x08, 0x00, 0x00})
	}))
	defer srv.Close()

	client := NewRemoteAPIClient("test-key", nil, nil, nil)
	client.baseURL = srv.URL

	_, err := client.makeRequestOnce("GET", "/v1/anything", nil)
	require.Error(t, err)

	var rateErr *RateLimitError
	assert.ErrorAs(t, err, &rateErr, "gzip decode failure must not mask 429 status")
}

// TestMakeRequestOnce_GzippedErrorBody_PreservesStatus ensures a 4xx with a
// malformed gzip body still surfaces a status-aware error rather than the
// generic "decoding response" error.
func TestMakeRequestOnce_GzippedErrorBody_PreservesStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		// Magic prefix + truncated stream.
		_, _ = w.Write([]byte{0x1f, 0x8b, 0x08, 0x00, 0x00})
	}))
	defer srv.Close()

	client := NewRemoteAPIClient("test-key", nil, nil, nil)
	client.baseURL = srv.URL

	_, err := client.makeRequestOnce("GET", "/v1/anything", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403", "decompression failure must not hide HTTP status")
}

func TestDiscoverSites_PlainJSONBody(t *testing.T) {
	t.Parallel()

	resp := SitesResponse{
		Data: []RemoteSite{{ID: "s1", InternalReference: "default", Name: "Default"}},
	}

	payload, err := json.Marshal(resp)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	client := NewRemoteAPIClient("test-key", nil, nil, nil)
	client.baseURL = srv.URL

	sites, err := client.DiscoverSites("console-1")
	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, "s1", sites[0].ID)
	assert.Equal(t, "default", sites[0].InternalReference)
	assert.Equal(t, "Default", sites[0].Name)
}

// TestDiscoverSites_OpenAPISiteOverviewPage locks in the Network Integration
// "Site overview page" shape (v10.4.57): id is the UUID, internalReference is
// the legacy site name, name is the display name. Pagination fields on the
// envelope are ignored. See unpoller/unpoller#986.
func TestDiscoverSites_OpenAPISiteOverviewPage(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
		"count": 2,
		"data": [
			{
				"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				"internalReference": "default",
				"name": "Default"
			},
			{
				"id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
				"internalReference": "abc1def2",
				"name": "Office"
			}
		],
		"limit": 25,
		"offset": 0,
		"totalCount": 2
	}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	client := NewRemoteAPIClient("test-key", nil, nil, nil)
	client.baseURL = srv.URL

	sites, err := client.DiscoverSites("console-1")
	require.NoError(t, err)
	require.Len(t, sites, 2)

	assert.Equal(t, "3fa85f64-5717-4562-b3fc-2c963f66afa6", sites[0].ID)
	assert.Equal(t, "default", sites[0].InternalReference)
	assert.Equal(t, "Default", sites[0].Name)

	assert.Equal(t, "7c9e6679-7425-40de-944b-e07fc1f90ae7", sites[1].ID)
	assert.Equal(t, "abc1def2", sites[1].InternalReference)
	assert.Equal(t, "Office", sites[1].Name)
}
