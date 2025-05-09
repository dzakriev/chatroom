package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func Produce() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
	})

	message := kafka.Message{
		Key:   []byte("Key1"),
		Value: []byte("Hello Kafka!"),
	}

	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent successfully!")

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}
}
