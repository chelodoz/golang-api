package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	HTTPServerAddress    string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string `mapstructure:"GRPC_SERVER_ADDRESS"`
	AccessTokenDuration  string `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration string `mapstructure:"REFRESH_TOKEN_DURATION"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBDriver             string `mapstructure:"DB_DRIVER"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DBName               string `mapstructure:"DB_NAME"`
	DBPort               string `mapstructure:"DB_PORT"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
