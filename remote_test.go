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
}

// TestDiscoverSites_GzipBodyWithoutContentEncoding reproduces unpoller/unpoller#997:
// the Site Manager proxy returns a gzipped body but does not set
// Content-Encoding: gzip, so Go's transport hands raw compressed bytes to the
// client. The library must detect and decompress the body itself.
func TestDiscoverSites_GzipBodyWithoutContentEncoding(t *testing.T) {
	t.Parallel()

	resp := SitesResponse{
		Data: []RemoteSite{
			{ID: "site-1", Name: "default", Description: "Default"},
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
	assert.Equal(t, "default", sites[0].Name)
}

func TestDiscoverSites_PlainJSONBody(t *testing.T) {
	t.Parallel()

	resp := SitesResponse{
		Data: []RemoteSite{{ID: "s1", Name: "default"}},
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
	assert.Equal(t, "default", sites[0].Name)
}
