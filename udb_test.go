package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestUDB(test *testing.T) {
	udb := unifi.UDB{}
	err := gofakeit.Struct(&udb)
	require.NoError(test, err)
	require.NotEmpty(test, udb.Name)
}
