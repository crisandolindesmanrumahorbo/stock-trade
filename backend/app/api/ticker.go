package api

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"github.com/stock-trade-app/core"
)

func GetAllTickers(c *fiber.Ctx) error {
	c.Status(200).JSON(core.GetAllTickers())
	return nil
}

func ListenTicker(conn *websocket.Conn) {
	currTicker := conn.Params("ticker")
	log.Println("Request Ticker: ", currTicker)
	if !core.IsTickerAllowed(currTicker) {
		conn.WriteMessage(websocket.CloseUnsupportedData, []byte("Ticker is not allowed"))
		log.Println("Ticker not allowed: ", currTicker)
		conn.Close()
		return
	}
	topic := "trades-" + strings.ToLower(currTicker)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{core.KAFKA_HOST + ":" + core.KAFKA_PORT},
		Topic:   topic,
	})
	reader.SetOffset(-1)

	conn.SetCloseHandler(func(code int, text string) error {
		reader.Close()
		log.Printf("Received connection close request. Closing connection kafka ...")
		return nil
	})

	defer conn.Close()
	defer reader.Close()
	go func() {
		code, wsMessage, err := conn.NextReader()
		if err != nil {
			log.Printf("err %s code %d", err.Error(), code)
			log.Println("Error reading last message from WS connection. Exiting ...")
			conn.Close()
			return
		}
		fmt.Printf("CODE: %d	 MESSAGE: %s", code, wsMessage)
	}()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error: ", err.Error())
			return
		}
		fmt.Println("Reading ...", string(message.Value))
		conn.WriteMessage(websocket.TextMessage, message.Value)
	}
}
