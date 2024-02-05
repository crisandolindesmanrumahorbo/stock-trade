package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/stock-trade-producer/trades"
)

var topics []string
var HOST string
var PORT string

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Print("Failed to load env variable")
	}
	t := os.Getenv("TICKERS")
	HOST = os.Getenv("KAFKA_HOST")
	PORT = os.Getenv("KAFKA_PORT")
	topics = strings.Split(t, ",")
}

func init() {
	loadEnv()
	trades.LoadHostAndPort(HOST, PORT)
}

func main() {
	trades.SubscribeAndListen(topics)
}
