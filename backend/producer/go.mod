module github.com/stock-trade-producer

go 1.21.5

require (
	github.com/gorilla/websocket v1.5.1 // to fetch from binance websocket api
	github.com/joho/godotenv v1.5.1 // to get tickers from env file
	github.com/segmentio/kafka-go v0.4.42 //sdk to connect with kafka
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.17.0 // indirect
)
