package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("missing file", func(t *testing.T) {
		_, err := loadConfig("/nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reading config")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		f, err := os.CreateTemp("", "unifi-config-*.json")
		require.NoError(t, err)

		defer os.Remove(f.Name())

		_, _ = f.WriteString("not json")
		require.NoError(t, f.Close())

		_, err = loadConfig(f.Name())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parsing config")
	})

	t.Run("valid config", func(t *testing.T) {
		f, err := os.CreateTemp("", "unifi-config-*.json")
		require.NoError(t, err)

		defer os.Remove(f.Name())

		_, _ = f.WriteString(`{"url":"https://controller:8443","user":"admin","pass":"secret"}`)
		require.NoError(t, f.Close())

		cfg, err := loadConfig(f.Name())
		require.NoError(t, err)
		assert.Equal(t, "https://controller:8443", cfg.URL)
		assert.Equal(t, "admin", cfg.User)
		assert.Equal(t, "secret", cfg.Pass)
	})
}

func TestGetEnvString(t *testing.T) {
	envName := "ABCXYZ_STRING"
	fallback := "Fallback"
	expected := "Expected"

	// Ensure env is clean before and after test.
	originalValue, isSet := os.LookupEnv(envName)
	os.Unsetenv(envName)

	defer func() {
		if isSet {
			os.Setenv(envName, originalValue)
		}
	}()

	// Test fallback when not set
	result := GetEnvString(envName, fallback)
	assert.Equal(t, fallback, result)

	// Test fallback with empty string
	os.Setenv(envName, "")
	result = GetEnvString(envName, fallback)
	assert.Equal(t, "", result)

	// Test with value set
	os.Setenv(envName, expected)
	result = GetEnvString(envName, fallback)
	assert.Equal(t, expected, result)
}

func TestGetEnvInt64(t *testing.T) {
	envName := "ABCXYZ_INT64"

	var fallback int64 = 123456

	var expected int64 = 654321

	originalValue, isSet := os.LookupEnv(envName)
	os.Unsetenv(envName)

	defer func() {
		if isSet {
			os.Setenv(envName, originalValue)
		}
	}()

	var result int64

	result = GetEnvInt64(envName, fallback)
	assert.Equal(t, fallback, result)

	os.Setenv(envName, "")
	result = GetEnvInt64(envName, fallback)
	assert.Equal(t, fallback, result)

	os.Setenv(envName, " ")
	result = GetEnvInt64(envName, fallback)
	assert.Equal(t, fallback, result)

	os.Setenv(envName, strconv.FormatInt(expected, 10))
	result = GetEnvInt64(envName, fallback)
	assert.Equal(t, expected, result)
}

func TestGetEnvInt(t *testing.T) {
	envName := "ABCXYZ_INT"

	var fallback = 123456

	var expected = 654321

	originalValue, isSet := os.LookupEnv(envName)
	os.Unsetenv(envName)

	defer func() {
		if isSet {
			os.Setenv(envName, originalValue)
		}
	}()

	var result int

	result = GetEnvInt(envName, fallback)
	assert.Equal(t, fallback, result)

	os.Setenv(envName, "")
	result = GetEnvInt(envName, fallback)
	assert.Equal(t, fallback, result)

	os.Setenv(envName, " ")
	result = GetEnvInt(envName, fallback)
	assert.Equal(t, fallback, result)

	os.Setenv(envName, strconv.Itoa(expected))
	result = GetEnvInt(envName, fallback)
	assert.Equal(t, expected, result)
}

func TestShowWithStructData(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}

	testData := []TestStruct{
		{ID: 1, Name: "First", Value: 10.5},
		{ID: 2, Name: "Second", Value: 20.5},
		{ID: 3, Name: "Third", Value: 30.5},
	}

	testCases := []struct {
		name         string
		prefix       string
		numResponses int
	}{
		{
			name:         "ShowResponse first item",
			prefix:       "StructItems",
			numResponses: 1,
		},
		{
			name:         "ShowResponse all items",
			prefix:       "AllStructItems",
			numResponses: 3,
		},
		{
			name:         "ShowResponse no items",
			prefix:       "NoStructItems",
			numResponses: 0,
		},
		{
			name:         "ShowResponse more than available",
			prefix:       "ExcessStructItems",
			numResponses: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				ShowResponse(tc.prefix, testData, tc.numResponses)
			})
		})
	}
}
