package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	AppPort = "app.port"
)

func LoadConfig() {
	getEnv()
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			log.Printf("Warning: Config file not found, using defaults/env")
		} else {
			log.Fatalf("Error: Failed to read config: %v", err)
		}
	}
	applyDefaults()
	if err := validateConfigFields(); err != nil {
		log.Fatalf("Fatal: config validation error: %s", err)
	}
}

func getEnv() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	var binds = map[string]string{
		AppPort: "PORT",
	}
	for k, v := range binds {
		_ = viper.BindEnv(k, v)
	}
}

func applyDefaults() {
	var defaults = map[string]any{AppPort: 8080} // Will be set if not present, overwrites above required/dependent

	for k, v := range defaults {
		if !viper.IsSet(k) {
			log.Printf("Warning: config field '%s' is missing or empty, defaulting to '%s'", k, v)
			viper.Set(k, v)
		}
	}
}

func validateConfigFields() error {
	var required = []string{AppPort} // Must be present and non-empty

	var missing []string
	for _, key := range required {
		if !viper.IsSet(key) {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields/values in config: %s", strings.Join(missing, ", "))
	}

	return nil
}
