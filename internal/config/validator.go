package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"abb_ia/internal/logger"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ConfigValidator validates configuration settings
type ConfigValidator struct {
	config *Config
	errors []ValidationError
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator(config *Config) *ConfigValidator {
	return &ConfigValidator{
		config: config,
		errors: make([]ValidationError, 0),
	}
}

// Validate performs all validation checks
func (v *ConfigValidator) Validate() error {
	v.validatePaths()
	v.validateLogSettings()
	v.validateNetworkSettings()
	v.validateProcessingSettings()

	if len(v.errors) > 0 {
		return v.formatErrors()
	}

	return nil
}

func (v *ConfigValidator) validatePaths() {
	// Validate output directory
	if v.config.OutputDir != "" {
		if !filepath.IsAbs(v.config.OutputDir) {
			v.addError("OutputDir", "must be an absolute path")
		}
		if _, err := os.Stat(v.config.OutputDir); os.IsNotExist(err) {
			// Try to create the directory
			if err := os.MkdirAll(v.config.OutputDir, 0755); err != nil {
				v.addError("OutputDir", "directory does not exist and cannot be created")
			}
		}
	}

	// Validate temporary directory
	if v.config.TmpDir != "" {
		if !filepath.IsAbs(v.config.TmpDir) {
			v.addError("TmpDir", "must be an absolute path")
		}
		if _, err := os.Stat(v.config.TmpDir); os.IsNotExist(err) {
			if err := os.MkdirAll(v.config.TmpDir, 0755); err != nil {
				v.addError("TmpDir", "directory does not exist and cannot be created")
			}
		}
	}
}

func (v *ConfigValidator) validateLogSettings() {
	// Validate log file
	if v.config.LogFileName != "" {
		dir := filepath.Dir(v.config.LogFileName)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				v.addError("LogFileName", "log directory cannot be created")
			}
		}
	}

	// Validate log level
	validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	logLevel := strings.ToUpper(v.config.LogLevel)
	isValidLevel := false
	for _, level := range validLevels {
		if logLevel == level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		v.addError("LogLevel", fmt.Sprintf("must be one of: %s", strings.Join(validLevels, ", ")))
	}
}

func (v *ConfigValidator) validateNetworkSettings() {
	// Validate concurrent downloaders
	if v.config.ConcurrentDownloaders < 1 {
		v.addError("ConcurrentDownloaders", "must be at least 1")
	}
	if v.config.ConcurrentDownloaders > 10 {
		logger.Warn(fmt.Sprintf("High number of concurrent downloaders may impact performance: %d", v.config.ConcurrentDownloaders))
	}
}

func (v *ConfigValidator) validateProcessingSettings() {
	// Validate concurrent encoders
	if v.config.ConcurrentEncoders < 1 {
		v.addError("ConcurrentEncoders", "must be at least 1")
	}
	if v.config.ConcurrentEncoders > 10 {
		logger.Warn(fmt.Sprintf("High number of concurrent encoders may impact performance: %d", v.config.ConcurrentEncoders))
	}

	// Validate bit rate
	if v.config.BitRateKbs < 32 || v.config.BitRateKbs > 320 {
		v.addError("BitRateKbs", fmt.Sprintf("must be between 32 and 320, got %d", v.config.BitRateKbs))
	}

	// Validate sample rate
	if v.config.SampleRateHz < 8000 || v.config.SampleRateHz > 48000 {
		v.addError("SampleRateHz", fmt.Sprintf("must be between 8000 and 48000, got %d", v.config.SampleRateHz))
	}

	// Validate max file size
	if v.config.MaxFileSizeMb < 1 || v.config.MaxFileSizeMb > 2000 {
		v.addError("MaxFileSizeMb", fmt.Sprintf("must be between 1 and 2000, got %d", v.config.MaxFileSizeMb))
	}
}

func (v *ConfigValidator) addError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v *ConfigValidator) formatErrors() error {
	var errMsgs []string
	for _, err := range v.errors {
		errMsgs = append(errMsgs, err.Error())
	}
	return fmt.Errorf("configuration validation failed:\n- %s", strings.Join(errMsgs, "\n- "))
}
