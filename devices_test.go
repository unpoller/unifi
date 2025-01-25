package unifi_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
	"testing"
)

func TestDevices(test *testing.T) {
	d := unifi.Devices{}
	err := gofakeit.Struct(&d)
	require.NoError(test, err)
	require.NotEmpty(test, d.UBBs)
	require.NotEmpty(test, d.UCIs)
	require.NotEmpty(test, d.PDUs)
	require.NotEmpty(test, d.UAPs)
	require.NotEmpty(test, d.USGs)
	require.NotEmpty(test, d.UDMs)
	require.NotEmpty(test, d.UXGs)
	require.NotEmpty(test, d.USWs)
}
