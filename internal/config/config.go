package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	JWTSecret   string
	DBHOST      string
	DBPORT      string
	DBUSER      string
	DBPASS      string
	DBNAME      string
	KafkaBroker string
	KafkaTopic  string
}

func Load() (*Config, error) {

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Port:        viper.GetString("PORT"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		DBHOST:      viper.GetString("DB_HOST"),
		DBPORT:      viper.GetString("DB_PORT"),
		DBUSER:      viper.GetString("DB_USER"),
		DBPASS:      viper.GetString("DB_PASS"),
		DBNAME:      viper.GetString("DB_NAME"),
		KafkaBroker: viper.GetString("KAFKA_BROKER"),
		KafkaTopic:  viper.GetString("KAFKA_TOPIC"),
	}

	return &cfg, nil
}
