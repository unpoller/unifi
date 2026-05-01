package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestLAG(t *testing.T) {
	t.Parallel()

	var l unifi.LAG

	require.NoError(t, gofakeit.Struct(&l))
}

func TestLAGMember(t *testing.T) {
	t.Parallel()

	var m unifi.LAGMember

	require.NoError(t, gofakeit.Struct(&m))
}

func TestLAGMetadata(t *testing.T) {
	t.Parallel()

	var m unifi.LAGMetadata

	require.NoError(t, gofakeit.Struct(&m))
}

func TestMCLAGDomain(t *testing.T) {
	t.Parallel()

	var d unifi.MCLAGDomain

	require.NoError(t, gofakeit.Struct(&d))
}

func TestMCLAGPeer(t *testing.T) {
	t.Parallel()

	var p unifi.MCLAGPeer

	require.NoError(t, gofakeit.Struct(&p))
}

func TestMCLAGMetadata(t *testing.T) {
	t.Parallel()

	var m unifi.MCLAGMetadata

	require.NoError(t, gofakeit.Struct(&m))
}

func TestSwitchStack(t *testing.T) {
	t.Parallel()

	var s unifi.SwitchStack

	require.NoError(t, gofakeit.Struct(&s))
}

func TestSwitchStackMetadata(t *testing.T) {
	t.Parallel()

	var m unifi.SwitchStackMetadata

	require.NoError(t, gofakeit.Struct(&m))
}
