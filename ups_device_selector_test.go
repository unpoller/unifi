package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestUPSDeviceSelectorStruct(t *testing.T) {
	t.Parallel()

	var d unifi.UPSDeviceSelector

	err := gofakeit.Struct(&d)
	require.NoError(t, err)
	require.NotEmpty(t, d.ID)
	require.NotEmpty(t, d.MAC)
}
