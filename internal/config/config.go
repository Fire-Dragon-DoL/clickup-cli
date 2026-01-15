package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SpaceID       string `mapstructure:"space_id"`
	OutputFormat  string `mapstructure:"output_format"`
	StrictResolve bool   `mapstructure:"strict_resolve"`
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetDefault("output_format", "text")

	v.BindEnv("space_id", "CLICKUP_SPACE_ID")
	v.BindEnv("output_format", "CLICKUP_OUTPUT_FORMAT")
	v.BindEnv("strict_resolve", "CLICKUP_STRICT_RESOLVE")

	return v
}

func Load() *Config {
	v := newViper()

	cfg := &Config{}
	v.Unmarshal(cfg)
	return cfg
}

func LoadFromFile(path string) *Config {
	v := newViper()

	v.SetConfigFile(path)
	v.ReadInConfig()

	cfg := &Config{}
	v.Unmarshal(cfg)
	return cfg
}

func (c *Config) ApplyCLIOverrides(spaceID, outputFormat string, strictResolve bool) {
	if spaceID != "" {
		c.SpaceID = spaceID
	}
	if outputFormat != "" {
		c.OutputFormat = outputFormat
	}
	if strictResolve {
		c.StrictResolve = true
	}
}
