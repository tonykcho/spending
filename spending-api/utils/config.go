package utils

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type ConnectionStrings struct {
	PostgresConnection string `json:"PostgresConnection"`
	DatabaseConnection string `json:"DatabaseConnection"`
	PaddleOcrHost      string `json:"PaddleOcrHost"`
	OllamaHost         string `json:"OllamaHost"`
	BasicAuthUser      string
	BasicAuthPassword  string
	Jaeger             string `json:"Jaeger"`
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

func GetPaddleOcrHost() string {
	if AppConfig.ConnectionStrings.PaddleOcrHost == "" {
		loadConfig()
	}
	return AppConfig.ConnectionStrings.PaddleOcrHost
}

func GetOllamaHost() string {
	if AppConfig.ConnectionStrings.OllamaHost == "" {
		loadConfig()
	}
	return AppConfig.ConnectionStrings.OllamaHost
}

func GetJaeger() string {
	if AppConfig.ConnectionStrings.Jaeger == "" {
		loadConfig()
	}
	return AppConfig.ConnectionStrings.Jaeger
}

func GetBasicAuthUser() string {
	if AppConfig.ConnectionStrings.BasicAuthUser == "" {
		loadAuthFromEnv()
	}
	return AppConfig.ConnectionStrings.BasicAuthUser
}

func GetBasicAuthPassword() string {
	if AppConfig.ConnectionStrings.BasicAuthPassword == "" {
		loadAuthFromEnv()
	}
	return AppConfig.ConnectionStrings.BasicAuthPassword
}

func loadAuthFromEnv() {
	user := os.Getenv("BASIC_AUTH_USER")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	AppConfig.ConnectionStrings.BasicAuthUser = user
	AppConfig.ConnectionStrings.BasicAuthPassword = password
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
