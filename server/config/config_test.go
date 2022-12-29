package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicConfiguration(t *testing.T) {
	expectedConfig := Config{
		ClickhouseConnUrl: "http://localhost:9000?username=admin&password=admin",
	}

	actual, err := FromFile("./testdata/basic_config.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, actual)
}
