package utils

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type ConnectionStrings struct {
	PostgresConnection string `json:"PostgresConnection"`
	DatabaseConnection string `json:"DatabaseConnection"`
}

type Config struct {
	ConnectionStrings ConnectionStrings `json:"ConnectionStrings"`
}

var AppConfig Config

func GetPostgresConnection() string {
	if AppConfig.ConnectionStrings.PostgresConnection == "" {
		loadConfig()
	}
	return AppConfig.ConnectionStrings.PostgresConnection
}

func GetDatabaseConnection() string {
	if AppConfig.ConnectionStrings.DatabaseConnection == "" {
		loadConfig()
	}
	return AppConfig.ConnectionStrings.DatabaseConnection
}

func loadConfig() {
	env := os.Getenv("APP_ENV")
	configFilePath := "config.json"
	switch env {
	case "dev":
		log.Info().Msg("Current env is dev, Loading development configuration")
		configFilePath = "config.dev.json"
	case "test":
		log.Info().Msg("Current env is test, Loading test configuration")
		configFilePath = "config.test.json"
	}

	file, err := os.Open(configFilePath)
	CheckError(err)
	defer file.Close()

	err = json.NewDecoder(file).Decode(&AppConfig)
	CheckError(err)
}
