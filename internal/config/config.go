package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig holds the configuration values for the application.
type AppConfig struct {
	SlackAppToken string `mapstructure:"slack_app_token"`
	SlackBotToken string `mapstructure:"slack_bot_token"`
	Environment   string `mapstructure:"environment"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(configPaths []string, configName string) (*AppConfig, error) {
	var config AppConfig

	viper.SetConfigName(configName) // Name of config file (without extension)
	viper.SetConfigType("toml")     // Set the type of the configuration file explicitly
	for _, path := range configPaths {
		viper.AddConfigPath(path) // Paths to look for the config file in
	}

	viper.AutomaticEnv() // Read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return &config, nil
}
