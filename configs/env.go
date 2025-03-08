package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PG_USER     string
	PG_PASSWORD string
	PG_HOST     string
	PG_DB       string
	PG_PORT     string
	PORT        string
	PRIV        []byte
	PUB         []byte
)

func InitEnv() {
	err := godotenv.Load("/app/.env")

	if err != nil {
		log.Fatal("Unable to load environment file")
	}

	PG_USER = getEnvOrDefault("POSTGRES_USER", "gorm")
	PG_PASSWORD = getEnvOrDefault("POSTGRES_PASSWORD", "gorm")
	PG_HOST = getEnvOrDefault("POSTGRES_HOST", "gorm")
	PG_DB = getEnvOrDefault("POSTGRES_DB", "gorm")
	PG_PORT = getEnvOrDefault("POSTGRES_PORT", "5432")
	PORT = getEnvOrDefault("PORT", "8000")
}

func LoadKeys() {
	var err error
	privFilePath := getEnvOrDefault("PRIVATE_KEY_FILE", "")
	PRIV, err = os.ReadFile(privFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	pubFilePath := getEnvOrDefault("PUBLIC_KEY_FILE", "")
	PUB, err = os.ReadFile(pubFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func getEnvOrDefault(s ...string) string {
	if len(s) <= 0 {
		return ""
	} else if val := os.Getenv(s[0]); len(s) >= 2 && val != "" {
		return val
	} else {
		return s[1]
	}
}
