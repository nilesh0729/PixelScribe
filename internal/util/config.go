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

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// Ignore file not found error - use environment variables instead
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// It's a different error, return it
			return
		}
		// File not found is OK - clear the error
		err = nil
	}
	
	err = viper.Unmarshal(&config)
	return
}

