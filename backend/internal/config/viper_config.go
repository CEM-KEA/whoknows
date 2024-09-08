package config

import (
	"strings"

	"github.com/spf13/viper"
)

// ViperConfig defines the interface for the viper configuration - without this interface, we cannot mock the viper configuration
type ViperConfig interface {
	SetConfigFile(file string)
	ReadInConfig() error
	SetConfigName(name string)
	SetConfigType(typ string)
	AddConfigPath(path string)
	AutomaticEnv()
	SetEnvPrefix(prefix string)
	SetEnvKeyReplacer(r *strings.Replacer)
	SetDefault(key string, value interface{})
	Unmarshal(rawVal interface{}) error
}

// ViperConfigWrapper is a wrapper around the viper configuration
type ViperConfigWrapper struct{}

// SetConfigFile sets the configuration file
func (v *ViperConfigWrapper) SetConfigFile(file string) {
	viper.SetConfigFile(file)
}

// ReadInConfig reads the configuration file
func (v *ViperConfigWrapper) ReadInConfig() error {
	return viper.ReadInConfig()
}

// SetConfigName sets the configuration name
func (v *ViperConfigWrapper) SetConfigName(name string) {
	viper.SetConfigName(name)
}

// SetConfigType sets the configuration type
func (v *ViperConfigWrapper) SetConfigType(typ string) {
	viper.SetConfigType(typ)
}

// AddConfigPath adds a configuration path
func (v *ViperConfigWrapper) AddConfigPath(path string) {
	viper.AddConfigPath(path)
}

// AutomaticEnv enables automatic environment variable reading
func (v *ViperConfigWrapper) AutomaticEnv() {
	viper.AutomaticEnv()
}

// SetEnvPrefix sets the environment prefix
func (v *ViperConfigWrapper) SetEnvPrefix(prefix string) {
	viper.SetEnvPrefix(prefix)
}

// SetEnvKeyReplacer sets the environment key replacer
func (v *ViperConfigWrapper) SetEnvKeyReplacer(r *strings.Replacer) {
	viper.SetEnvKeyReplacer(r)
}

// SetDefault sets the default value for a key
func (v *ViperConfigWrapper) SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

// Unmarshal unmarshals the configuration into a struct
func (v *ViperConfigWrapper) Unmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}
