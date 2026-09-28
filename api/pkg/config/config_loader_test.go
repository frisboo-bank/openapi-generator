package config_test

import (
	"reflect"
	"testing"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/config"
	"frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type colorType int8

const (
	colorUnknown colorType = iota
	colorRed
)

type testAppConfig struct {
	Name    string        `mapstructure:"name"`
	Timeout time.Duration `mapstructure:"timeout"`
	Color   colorType     `mapstructure:"color"`
}

type testConfig struct {
	App testAppConfig `mapstructure:"app"`
}

func colorDecodeHook() mapstructure.DecodeHookFunc {
	return func(_, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeFor[colorType]() {
			return data, nil
		}
		if s, ok := data.(string); ok && s == "red" {
			return colorRed, nil
		}
		return data, nil
	}
}

func createConfigLoader(t *testing.T) contracts.ConfigLoader {
	t.Helper()
	loader, err := config.NewConfigLoader(config.ConfigLoaderOptions{
		ConfigName: "application",
		ConfigPath: "testdata",
	}, viper.New())
	require.NoError(t, err)
	loader.RegisterDecodeHookFunc(colorDecodeHook())
	return loader
}

func TestConfigLoader(t *testing.T) {
	loader := createConfigLoader(t)

	var cfg testConfig
	err := loader.Load(environmentenum.Environments.DEVELOPMENT, &cfg)
	require.NoError(t, err)

	expected := testConfig{
		App: testAppConfig{
			Name:    "test-service",
			Timeout: 10 * time.Second,
			Color:   colorRed,
		},
	}

	assert.Equal(t, expected, cfg)
}

func TestLoadKey(t *testing.T) {
	loader := createConfigLoader(t)

	var appConfig testAppConfig
	err := loader.LoadKey(environmentenum.Environments.PRODUCTION, &appConfig, "app")
	require.NoError(t, err)

	expected := testAppConfig{
		Name:    "test-service (prod)",
		Timeout: 5 * time.Second,
		Color:   colorRed,
	}

	assert.Equal(t, expected, appConfig)
}

func TestLoadKeyLeaf(t *testing.T) {
	loader := createConfigLoader(t)

	var color colorType
	err := loader.LoadKey(environmentenum.Environments.DEVELOPMENT, &color, "app.color")
	require.NoError(t, err)

	assert.Equal(t, colorRed, color)
}

func TestHasKey(t *testing.T) {
	loader := createConfigLoader(t)

	ok, err := loader.HasKey(environmentenum.Environments.DEVELOPMENT, "app")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = loader.HasKey(environmentenum.Environments.DEVELOPMENT, "app.color")
	require.NoError(t, err)
	assert.True(t, ok)

	_, err = loader.HasKey(environmentenum.Environments.DEVELOPMENT, "")
	assert.Error(t, err)
}

func TestNewConfigLoader_NilViper(t *testing.T) {
	loader, err := config.NewConfigLoader(config.ConfigLoaderOptions{}, nil)
	assert.Nil(t, loader)
	assert.Error(t, err)
}

func TestMissingKey(t *testing.T) {
	loader := createConfigLoader(t)

	ok, err := loader.HasKey(environmentenum.Environments.DEVELOPMENT, "missing")
	assert.False(t, ok)
	assert.NoError(t, err)

	var cfg testConfig
	err = loader.LoadKey(environmentenum.Environments.DEVELOPMENT, &cfg, "missing")
	assert.Error(t, err)
}
