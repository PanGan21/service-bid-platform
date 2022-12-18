package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Publisher interface {
	Publish(topic string, message Message) error
	Close() error
}

type kafkaPublisher struct {
	publisher *kafka.Writer
	retries   int
}

func NewPublisher(url string, retries int) Publisher {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(url),
		AllowAutoTopicCreation: false, // topics defined in environment variables in order to specify the partitions
	}

	return &kafkaPublisher{
		publisher: w,
		retries:   retries,
	}
}

func (k kafkaPublisher) Publish(topic string, msg Message) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = k.publisher.WriteMessages(ctx,
			kafka.Message{Topic: topic, Value: msgJson})
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err == nil {
			break
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}
	}

	if err != nil {
		fmt.Println("failed to write messages:", err)
		return err
	}

	fmt.Printf("PUBLISHER: Topic: %s - Payload: %v\n", topic, msg.Payload)

	return nil
}

func (k kafkaPublisher) Close() error {
	return k.publisher.Close()
}
