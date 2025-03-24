package initializers

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AuthDomain       string `mapstructure:"AUTH0_DOMAIN"`
	AuthClientID     string `mapstructure:"AUTH0_CLIENT_ID"`
	AuthClientSecret string `mapstructure:"AUTH0_CLIENT_SECRET"`
	AuthAudience     string `mapstructure:"AUTH0_AUDIENCE"`
	AuthRedirectURL  string `mapstructure:"AUTH0_CALLBACK_URL"`
}

var config Config

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func GetConfig() *Config {
	return &config
}
