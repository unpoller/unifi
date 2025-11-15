package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestClientUsageByApp(test *testing.T) {
	t := unifi.ClientUsageByApp{}
	err := gofakeit.Struct(&t)
	require.NoError(test, err)
	require.NotEmpty(test, t.UsageByApp)
}

func TestUsageByCountry(test *testing.T) {
	t := unifi.UsageByCountry{}
	err := gofakeit.Struct(&t)
	require.NoError(test, err)
}
