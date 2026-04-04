package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	LogLevel       = "app.log.level"           // string ("DEBUG", "INFO", "WARN", "ERROR")
	LogFormat      = "app.log.log_format"      // string ("text" or "json")
	LogToConsole   = "app.log.log2console"     // bool
	LogToFile      = "app.log.log2file"        // bool
	LogFilePath    = "app.log.file_path"       // string (path)
	LogFileMode    = "app.log.file_mode"       // string ("append", "overwrite", "rotate")
	LogFilesFolder = "app.log.old_logs_folder" // string (path)

	GinReleaseMode           = "app.web.gin_release_mode" // bool
	AppHost                  = "app.web.host"             // string
	AppPort                  = "app.web.port"             // int
	WebServerShutdownTimeout = "app.web.shutdown_timeout" // time.Duration
)

func LoadConfig() {
	getEnv()
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		var vNotFound viper.ConfigFileNotFoundError
		var osNotFound *fs.PathError
		if errors.As(err, &vNotFound) || errors.As(err, &osNotFound) {
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
		/* Logger */ LogLevel: "LOG_LEVEL", LogFormat: "LOG_FORMAT", LogToConsole: "LOG_TO_CONSOLE", LogToFile: "LOG_TO_FILE", LogFilePath: "LOG_FILE_PATH", LogFileMode: "LOG_FILE_MODE", LogFilesFolder: "LOG_FILES_FOLDER",
		/* Web Server */ GinReleaseMode: "GIN_RELEASE_MODE", AppHost: "APP_HOST", AppPort: "APP_PORT", WebServerShutdownTimeout: "WEB_SERVER_SHUTDOWN_TIMEOUT",
	}
	for k, v := range binds {
		_ = viper.BindEnv(k, v)
	}
}

func applyDefaults() {
	var defaults = map[string]any{ // Will be set if not present, overwrites above required/dependent
		/* Logger */ LogLevel: "INFO", LogFormat: "text", LogToConsole: true, LogToFile: true, LogFilePath: "application.log", LogFileMode: "append",
		/* Web Server */ GinReleaseMode: true, AppHost: "0.0.0.0", AppPort: 8080, WebServerShutdownTimeout: "5s",
	}
	for k, v := range defaults {
		if !viper.IsSet(k) || strings.TrimSpace(viper.GetString(k)) == "" {
			log.Printf("Warning: config field '%s' is missing or empty, defaulting to '%v'\n", k, v)
			viper.Set(k, v)
		}
	}
}

func validateConfigFields() error {
	var requiredFields = []string{ // Must be present and non-empty
		/* Web server */ AppHost, AppPort, WebServerShutdownTimeout,
	}
	var dependentFields = map[string][]string{ // E.g. if A=true ==> must be non-empty B
		LogToFile: {LogFilePath},
	}
	var possibleValues = map[string][]string{ // If key is present, value must be one of these values
		LogLevel:    {"DEBUG", "INFO", "WARN", "ERROR"},
		LogFormat:   {"text", "json"},
		LogFileMode: {"append", "overwrite", "rotate"},
	}
	var durationValues = []string{ // Will be checked if set duration is >= 0
		WebServerShutdownTimeout,
	}
	var conditionalRequired = map[string]map[string][]string{
		LogFileMode: {
			"rotate": {LogFilesFolder},
		},
	}

	var missing []string
	for _, key := range requiredFields {
		if isEmptyValue(key) {
			missing = append(missing, key)
		}
	}
	for triggerKey, requiredKeys := range dependentFields {
		if viper.GetBool(triggerKey) {
			for _, key := range requiredKeys {
				if isEmptyValue(key) {
					missing = append(missing, fmt.Sprintf("%s (required when %s=true)", key, triggerKey))
				}
			}
		}
	}
	for key, cases := range conditionalRequired {
		val := getValue(key)
		required, ok := cases[val]
		if !ok {
			continue
		}
		for _, field := range required {
			if isEmptyValue(field) {
				missing = append(missing, fmt.Sprintf("%s (required when %s=%s)", field, key, val))
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required fields/values in config: %s", strings.Join(missing, ", "))
	}

	var invalid []string
	for key, allowed := range possibleValues {
		val := getValue(key)
		if !viper.IsSet(key) || val == "" {
			continue
		}
		found := false
		for _, a := range allowed {
			if val == a {
				found = true
				break
			}
		}
		if !found {
			invalid = append(invalid, fmt.Sprintf("'%s' for '%s' (must be one of [%s])", val, key, strings.Join(allowed, ", ")))
		}
	}
	for _, key := range durationValues {
		if viper.GetDuration(key) <= 0 {
			invalid = append(invalid, fmt.Sprintf("%s (duration must be >0, got '%s')", key, viper.GetString(key)))
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("invalid config values: %s", strings.Join(invalid, ", "))
	}

	return nil
}

func getValue(key string) string {
	return strings.TrimSpace(viper.GetString(key))
}

func isEmptyValue(key string) bool {
	return getValue(key) == ""
}
