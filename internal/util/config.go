package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress      string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey  string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	OpenAIKey           string        `mapstructure:"OPENAI_API_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)       
	viper.AddConfigPath(".")        
	viper.AddConfigPath("../")      
	viper.AddConfigPath("../..")    
	viper.AddConfigPath("../../..") 

	viper.SetConfigName(".env")  // Look for .env (without extension)
	viper.SetConfigType("env")
	
	// Enable automatic env variable reading
	viper.AutomaticEnv()
	
	// Explicitly bind environment variables to config keys
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("OPENAI_API_KEY")

	// Try to read config file, but don't fail if it doesn't exist
	_ = viper.ReadInConfig()  // Ignore all errors from file reading
	
	// Unmarshal will use env vars if no file was found
	err = viper.Unmarshal(&config)
	return
}

