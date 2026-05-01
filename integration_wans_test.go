package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestIntegrationWAN(t *testing.T) {
	t.Parallel()

	var w unifi.IntegrationWAN

	err := gofakeit.Struct(&w)
	require.NoError(t, err)
	require.NotEmpty(t, w.ID)
	require.NotEmpty(t, w.Name)
}
