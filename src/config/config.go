package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	// LevelDebug config value for logging level `debug`.
	LevelDebug = "debug"
	// LevelInfo config value for logging level `info`.
	LevelInfo = "info"
)

var levels = []string{LevelDebug, LevelInfo}

// Config is this API's runtime config.
type Config struct {
	Port     int    `envconfig:"port" required:"true"`
	LogLevel string `envconfig:"log_level" required:"true"`

	BrokerConnectionString string `envconfig:"broker_connection_string" required:"true"`
}

// New will construct a config instance.
//
// NOTE: This is a failable constructor. If the configuration is not valid, this routine will panic.
func New() *Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	// Ensure log level is valid.
	if err := validateLogLevel(config.LogLevel); err != nil {
		panicWithArgs(err.Error())
	}

	return &config
}

/////////////////////
// Private Symbols //

// panicWithArgs will panic with the given arguments.
func panicWithArgs(errStr string) {
	panic(fmt.Sprintf("Invalid configuration. %s", errStr))
}

// validateLogLevel will validate that the `Mode` field has a valid value.
func validateLogLevel(level string) error {
	for _, validLevel := range levels {
		if level == validLevel {
			return nil
		}
	}
	return fmt.Errorf("Log level '%s' is invalid. Must be one of '%T'.", level, levels)
}
