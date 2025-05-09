package kafka

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker = "localhost:9092"
	kafkaTopic  = "message"
	groupID     = "message"
	redisAddr   = "localhost:6379"
	redisTTL    = 3600 * time.Second
)

func Consume() {
	go main()
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: groupID,
	})
	defer r.Close()

	ctx := context.Background()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message from Kafka:", err)
			continue
		}

		messageID := string(msg.Key)

		isDuplicate, err := isMessageDuplicate(ctx, rdb, messageID)
		if err != nil {
			log.Println("Error checking message ID in Redis:", err)
			continue
		}

		if isDuplicate {
			log.Printf("Duplicate message detected, skipping: %s\n", messageID)
			continue
		}

		log.Printf("Processing message: %s\n", messageID)

		err = storeMessageID(ctx, rdb, messageID)
		if err != nil {
			log.Println("Error storing message ID in Redis:", err)
		}
	}
}

func isMessageDuplicate(ctx context.Context, rdb *redis.Client, messageID string) (bool, error) {
	exists, err := rdb.Exists(ctx, messageID).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func storeMessageID(ctx context.Context, rdb *redis.Client, messageID string) error {
	err := rdb.Set(ctx, messageID, "processed", redisTTL).Err()
	return err
}
