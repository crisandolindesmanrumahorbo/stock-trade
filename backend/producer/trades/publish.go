package trades

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var (
	HOST string
	PORT string
)

func LoadHostAndPort(host string, port string) {
	HOST = host
	PORT = port
}

func Publish(message kafka.Message, topic string) error {
	messages := []kafka.Message{
		message,
	}
	w := kafka.Writer{
		Addr:                   kafka.TCP(HOST + ":" + PORT),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}
	defer w.Close()

	err := w.WriteMessages(context.Background(), messages...)
	if err != nil {
		log.Println("Error writing messages to Kafka: ", err.Error())
		return err
	}
	log.Println("Publish message to Kafka on topic: ", topic)
	return nil
}
