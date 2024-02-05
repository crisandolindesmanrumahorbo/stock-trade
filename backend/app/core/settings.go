package core

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	KAFKA_HOST string
	KAFKA_PORT string
	TICKERS    []string
)

func Load() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Failed to load env variable")
	}
	t := os.Getenv("TICKERS")
	TICKERS = strings.Split(t, ",")
	KAFKA_HOST = os.Getenv("KAFKA_HOST")
	KAFKA_PORT = os.Getenv("KAFKA_PORT")
	LoadTickers(TICKERS)
}
