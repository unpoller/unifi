package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestIntegrationNetwork(t *testing.T) {
	t.Parallel()

	var n unifi.IntegrationNetwork

	err := gofakeit.Struct(&n)
	require.NoError(t, err)
	require.NotEmpty(t, n.ID)
	require.NotEmpty(t, n.Name)
}
