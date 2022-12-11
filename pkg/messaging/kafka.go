package messaging

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kafkaMessagingService struct {
	broker string //"localhost:9092"
}

func NewMessagingService(broker string) *kafkaMessagingService {
	// conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
	return &kafkaMessagingService{
		broker: broker,
	}
}

func (k *kafkaMessagingService) Publish(ctx context.Context, topic string, msg Message) error {
	kafka.DialLeader(ctx, "tcp", k.broker, topic, msg.)

}

func Subscribe(topic string, groupId string, handler EventHandler) *Consumer {
	return nil
}

func newKafkaConfig(address string, topic string) (host, topic string, err error) {
	// host, err = conf.Get("KAFKA_HOST")
	// if err != nil {
	// 	return "", "", fmt.Errorf("conf.Get KAFKA_HOST %w", err)
	// }

	// topic, err = conf.Get("KAFKA_TOPIC")
	// if err != nil {
	// 	return "", "", fmt.Errorf("conf.Get KAFKA_TOPIC %w", err)
	// }

	if topic == "" {
		return "", "", internaldomain.NewErrorf(internaldomain.ErrorCodeInvalidArgument, "KAFKA_TOPIC is required")
	}

	return host, topic, nil
}
