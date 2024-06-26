package config

import "github.com/spf13/viper"

// Holds configuration information
type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUsername    string `mapstructure:"DB_USERNAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
}

// Load configuration from files
func Load() (*Config, error) {
	cfg := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
