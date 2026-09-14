package config

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/reflection"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
	"frisboo-bank/openapi-generator-service/pkg/utils"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var _ contracts.ConfigLoader = (*configLoader)(nil)

type ConfigLoaderOptions struct {
	ConfigName      string
	ConfigPath      string
	Debug           bool
	EnvKeyReplacer  map[string]string
	EnvPrefix       string
	DecodeHookFuncs []mapstructure.DecodeHookFunc
}

type configLoader struct {
	basePath string
	config   ConfigLoaderOptions
	loadErr  error
	loadOnce sync.Once
	mu       sync.RWMutex
	viper    *viper.Viper
}

func NewConfigLoader(
	cfg ConfigLoaderOptions,
	vi *viper.Viper,
) (contracts.ConfigLoader, error) {
	if vi == nil {
		return nil, fmt.Errorf("viper instance cannot be nil")
	}

	if cfg.Debug {
		vi.Debug()
	}

	if cfg.ConfigName != "" {
		vi.SetConfigName(cfg.ConfigName)
	}

	if cfg.ConfigPath != "" {
		vi.AddConfigPath(cfg.ConfigPath)
	}

	configPath := cfg.ConfigPath
	if !path.IsAbs(configPath) {
		root, err := utils.GetProjectRootWorkingDirectory()
		if err != nil {
			return nil, err
		}
		configPath = path.Join(root, configPath)
	}

	vi.AddConfigPath(configPath)
	vi.AddConfigPath(".")

	vi.AutomaticEnv()

	if cfg.EnvPrefix != "" {
		vi.SetEnvPrefix(cfg.EnvPrefix)
	}

	if cfg.EnvKeyReplacer != nil {
		replacer := make([]string, 0, len(cfg.EnvKeyReplacer)*2)
		for match, replace := range cfg.EnvKeyReplacer {
			replacer = append(replacer, match, replace)
		}
		vi.SetEnvKeyReplacer(strings.NewReplacer(replacer...))
	}

	return &configLoader{
		basePath: configPath,
		config:   cfg,
		viper:    vi,
	}, nil
}

func (c *configLoader) Load(env environmentEnum.Environment, target any) error {
	if !reflection.IsPointer(target) {
		return syserrors.Newf("target must be a pointer: got %s", reflection.GetKind(target))
	}
	if err := c.ensureLoaded(env); err != nil {
		return err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.unmarshal("", target)
}

func (c *configLoader) LoadKey(env environmentEnum.Environment, target any, key string) error {
	if !reflection.IsPointer(target) {
		return syserrors.Newf("target must be a pointer: got %s", reflection.GetKind(target))
	}
	if err := c.ensureLoaded(env); err != nil {
		return err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.unmarshal(key, target)
}

func (c *configLoader) HasKey(env environmentEnum.Environment, key string) (bool, error) {
	if key == "" {
		return false, syserrors.CantBeEmptyError("key")
	}
	if err := c.ensureLoaded(env); err != nil {
		return false, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.keyExists(key), nil
}

// Internal Helpers

func (c *configLoader) ensureLoaded(env environmentEnum.Environment) error {
	c.loadOnce.Do(func() {
		c.loadErr = c.doLoad(env)
	})
	return c.loadErr
}

func (c *configLoader) doLoad(env environmentEnum.Environment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.viper.SetConfigName(c.config.ConfigName)
	if err := c.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("base config %q: %w", c.config.ConfigName, err)
	}

	envConfigName := fmt.Sprintf("%s_%s", c.config.ConfigName, env)
	c.viper.SetConfigName(envConfigName)
	if err := c.viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("env config %q: %w", envConfigName, err)
		}
	}

	return nil
}

func (c *configLoader) unmarshal(key string, target any) error {
	hooks := []mapstructure.DecodeHookFunc{
		mapstructure.StringToTimeDurationHookFunc(),
	}
	hooks = append(hooks, c.config.DecodeHookFuncs...)

	opts := []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(hooks...)),
	}

	if key == "" {
		if err := c.viper.UnmarshalExact(target, opts...); err != nil {
			return syserrors.Wrap(err, "failed to unmarshal root")
		}
		return nil
	}

	sub := c.viper.Sub(key)
	if sub == nil || len(sub.AllSettings()) == 0 {
		return syserrors.Newf("required key %s not found", key)
	}

	if err := sub.UnmarshalExact(target, opts...); err != nil {
		return fmt.Errorf("failed to unmarshal key %s with error: %w", key, err)
	}

	return nil
}

func (c *configLoader) keyExists(key string) bool {
	if c.viper.IsSet(key) {
		return true
	}
	sub := c.viper.Sub(key)
	return sub != nil && len(sub.AllSettings()) > 0
}
