package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestVPNServer(t *testing.T) {
	t.Parallel()

	var s unifi.VPNServer

	err := gofakeit.Struct(&s)
	require.NoError(t, err)
	require.NotEmpty(t, s.ID)
	require.NotEmpty(t, s.Name)
}

func TestVPNServerMetadata(t *testing.T) {
	t.Parallel()

	var m unifi.VPNServerMetadata

	err := gofakeit.Struct(&m)
	require.NoError(t, err)
}

func TestSiteToSiteTunnel(t *testing.T) {
	t.Parallel()

	var tun unifi.SiteToSiteTunnel

	err := gofakeit.Struct(&tun)
	require.NoError(t, err)
	require.NotEmpty(t, tun.ID)
	require.NotEmpty(t, tun.Name)
}

func TestSiteToSiteTunnelMetadata(t *testing.T) {
	t.Parallel()

	var m unifi.SiteToSiteTunnelMetadata

	err := gofakeit.Struct(&m)
	require.NoError(t, err)
}
