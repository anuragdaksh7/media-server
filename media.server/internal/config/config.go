package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PORT          string `mapstructure:"PORT"`
	DbString      string `mapstructure:"DB_STRING"`
	JwtSecret     string `mapstructure:"JWT_SECRET"`
	Environment   string `mapstructure:"ENVIRONMENT"`
	EncryptKey    string `mapstructure:"ENCRYPTION_KEY"`
	RedisURL      string `mapstructure:"REDIS_URL"`
	StoragePath   string `mapstructure:"STORAGE_ROOT"`
	MaxUploadSize int64  `mapstructure:"MAX_UPLOAD_SIZE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	//viper.SetConfigName("app")
	//viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
