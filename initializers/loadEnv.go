package initializers

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AuthDomain   string `mapstructure:"AUTH0_DOMAIN"`
	AuthAudience string `mapstructure:"AUTH0_AUDIENCE"`
}

var config Config

func LoadConfig(path string) (err error) {

	// --- Direct Environment Variable Check ---
	log.Println("--- Checking ENV VARS directly via os.Getenv ---")
	pgUser := os.Getenv("POSTGRES_USER")
	pgHost := os.Getenv("POSTGRES_HOST")
	serverPort := os.Getenv("PORT")
	log.Printf("os.Getenv(\"POSTGRES_USER\"): [%s]", pgUser)
	log.Printf("os.Getenv(\"POSTGRES_HOST\"): [%s]", pgHost)
	log.Printf("os.Getenv(\"PORT\"): [%s]", serverPort)
	log.Println("--- End direct ENV VAR check ---")
	// --- End Direct Check ---

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	// Check if the error is specifically "Config File Not Found"
	// If it is, we IGNORE it because we want to rely on environment variables for railway deployment.
	// If it's any other error (e.g., parsing error), we return it.
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if we want to rely on ENV vars
			log.Println("Config file 'app.env' not found, relying on environment variables.")
			
			log.Println("Explicitly binding ENV VARS to Viper keys sinice AutomaticEnv doesn't work in railway?...")
			bindEnvErr := viper.BindEnv("POSTGRES_HOST")
			if bindEnvErr != nil { log.Printf("BindEnv POSTGRES_HOST error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("POSTGRES_USER")
			if bindEnvErr != nil { log.Printf("BindEnv POSTGRES_USER error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("POSTGRES_PASSWORD")
			if bindEnvErr != nil { log.Printf("BindEnv POSTGRES_PASSWORD error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("POSTGRES_DB")
			if bindEnvErr != nil { log.Printf("BindEnv POSTGRES_DB error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("POSTGRES_PORT")
			if bindEnvErr != nil { log.Printf("BindEnv POSTGRES_PORT error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("PORT")
			if bindEnvErr != nil { log.Printf("BindEnv PORT error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("CLIENT_ORIGIN")
			if bindEnvErr != nil { log.Printf("BindEnv CLIENT_ORIGIN error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("AUTH0_DOMAIN")
			if bindEnvErr != nil { log.Printf("BindEnv AUTH0_DOMAIN error: %v", bindEnvErr); return }
			bindEnvErr = viper.BindEnv("AUTH0_AUDIENCE")
			if bindEnvErr != nil { log.Printf("BindEnv AUTH0_AUDIENCE error: %v", bindEnvErr); return }
			
			err = nil // Clear the error so we don't return it
		} else {
			// Config file was found but another error occurred
			log.Printf("Error reading config file: %s", err)
			return 
		}
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Printf("!!! Viper Unmarshal Error: %v", err)
	} else {
		log.Printf("--- Config loaded via Viper ---")
		log.Printf("DBHost: [%s]", config.DBHost)
		log.Printf("DBUserName: [%s]", config.DBUserName)
		log.Printf("DBPassword: [%s]", "***")
		log.Printf("DBName: [%s]", config.DBName)
		log.Printf("DBPort: [%s]", config.DBPort)
		log.Printf("ServerPort: [%s]", config.ServerPort)
		log.Printf("ClientOrigin: [%s]", config.ClientOrigin)
		log.Printf("AuthDomain: [%s]", config.AuthDomain)
		log.Printf("AuthAudience: [%s]", config.AuthAudience)
		log.Printf("--- End Config ---")
	}

	return
}

func GetConfig() *Config {
	return &config
}
