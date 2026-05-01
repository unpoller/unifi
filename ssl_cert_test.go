package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestSSLCertificateStruct(t *testing.T) {
	t.Parallel()

	var c unifi.SSLCertificate

	err := gofakeit.Struct(&c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)
}

func TestSSLCertificateChainStruct(t *testing.T) {
	t.Parallel()

	var chain unifi.SSLCertificateChain

	err := gofakeit.Struct(&chain)
	require.NoError(t, err)
}

func TestNetworkExtendedFields(t *testing.T) {
	t.Parallel()

	var n unifi.Network

	err := gofakeit.Struct(&n)
	require.NoError(t, err)
}
