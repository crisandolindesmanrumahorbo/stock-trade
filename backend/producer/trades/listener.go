package trades

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

var conn *websocket.Conn

func getConnection() (*websocket.Conn, error) {
	if conn != nil {
		return conn, nil
	}
	u := url.URL{Scheme: "wss", Host: "ws.eodhistoricaldata.com", Path: "/ws/us", RawQuery: "api_token=demo"}
	log.Printf("connecting to %s", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	log.Printf("handshake succeed with status %d", resp.StatusCode)
	return c, nil
}

func unsubscribeOnClose(symbols string) error {
	message := struct {
		Action  string `json:"action"`
		Symbols string `json:"symbols"`
	}{
		Action:  "unsubscribe",
		Symbols: symbols,
	}
	m, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to JSON Encode subs message")
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Fatal("Failed to subscribe to topics ", err.Error())
		return err
	}
	return nil
}

func SubscribeAndListen(topics []string) error {
	conn, err := getConnection()
	if err != nil {
		log.Fatalf("Failed to get connection %s", err.Error())
		return err
	}
	conn.SetPongHandler(func(appData string) error {
		fmt.Println("Received pong: ", appData)
		pingFrame := []byte{1, 2, 3, 4, 5}
		err := conn.WriteMessage(websocket.PingMessage, pingFrame)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})

	symbols := fmt.Sprintln(strings.Join(topics[:], ","))
	message := struct {
		Action  string `json:"action"`
		Symbols string `json:"symbols"`
	}{
		Action:  "subscribe",
		Symbols: symbols,
	}
	m, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to JSON Encode subs message")
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Fatal("Failed to subscribe to topics ", err.Error())
		return err
	}
	log.Printf("succeed send message: %s", m)
	defer conn.Close()
	defer unsubscribeOnClose(symbols)

	for {
		_, payload, err := conn.ReadMessage()
		log.Printf("succeed received message: %s", payload)
		if err != nil {
			fmt.Println(err)
			return err
		}

		trade := Ticker{}

		err = json.Unmarshal(payload, &trade)
		if err != nil {
			fmt.Println(err)
			return err
		}
		log.Println(trade.Symbol, trade.Price, trade.Volume)

		// kafka send
		go convertAndPublishToKafka(trade)
	}
}

func convertAndPublishToKafka(t Ticker) {
	bytes, err := json.Marshal(t)
	if err != nil {
		log.Println("Error marshalling Ticker data ", err.Error())
	}
	message := kafka.Message{
		Key:   []byte(t.Symbol + "-" + strconv.Itoa(int(t.Time))),
		Value: bytes,
	}
	topic := "trades-" + strings.ToLower(t.Symbol)

	Publish(message, topic)
}
