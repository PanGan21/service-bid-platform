package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Subscriber interface {
	Subscribe(topic string, fn fnHandler)
}

type kafkaSubscriber struct {
	broker  string
	groupId string
}

type fnHandler func(payload interface{}) error

func NewSubscriber(url string, groupId string) Subscriber {
	return &kafkaSubscriber{broker: url, groupId: groupId}
}

func (k kafkaSubscriber) Subscribe(topic string, fn fnHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{k.broker},
		GroupID:        k.groupId,
		Topic:          topic,
		GroupBalancers: []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}},
		CommitInterval: time.Second,
	})

	defer r.Close()

	ctx := context.Background()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}

		var deserialized Message
		if err := json.Unmarshal(m.Value, &deserialized); err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("SUBSCRIBER: Topic: %s - Partition: %d - Offset: %d - Key: %s - Payload: %v\n", m.Topic, m.Partition, m.Offset, string(m.Key), deserialized.Payload)

		err = fn(deserialized.Payload)
		if err != nil {
			fmt.Println(err)
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
