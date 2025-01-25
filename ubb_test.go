package unifi_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
	"testing"
)

func TestUBB(test *testing.T) {
	ubb := unifi.UBB{}
	err := gofakeit.Struct(&ubb)
	require.NoError(test, err)
	require.NotEmpty(test, ubb.Name)
}
