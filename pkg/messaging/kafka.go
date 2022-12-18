package messaging

// var KAFKA_URL = os.Getenv("KAFKA_URL")

// type MessageBus interface {
// 	Publish(topic string, payload interface{}) error
// 	Subscribe(topic string, handler fnHandler) error
// }

// type kafkaMessageBus struct {
// 	Publisher  sarama.AsyncProducer
// 	Subscriber sarama.Consumer
// }

// type fnHandler func(payload []byte) error

// func NewMessageBus(kafkaUrl string) (MessageBus, error) {
// 	brokers := []string{kafkaUrl}
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true
// 	config.Producer.Retry.Max = 5
// 	config.Producer.Partitioner = sarama.NewRandomPartitioner
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Return.Successes = true
// 	config.Producer.Return.Errors = true

// 	subscriber, err := sarama.NewConsumer(brokers, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	publisher, err := sarama.NewAsyncProducer(brokers, config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &kafkaMessageBus{publisher, subscriber}, nil
// }

// func (bus *kafkaMessageBus) Publish(topic string, payload interface{}) error {
// 	msg, err := json.Marshal(payload)
// 	if err != nil {
// 		return err
// 	}

// 	message := &sarama.ProducerMessage{
// 		Topic:     topic,
// 		Partition: -1,
// 		Value:     sarama.StringEncoder(msg),
// 	}

// 	for {
// 		select {
// 		case bus.Publisher.Input() <- message:
// 		case <-bus.Publisher.Successes():
// 			return nil
// 		case err := <-bus.Publisher.Errors():
// 			return err
// 		}
// 	}
// }

// func (bus *kafkaMessageBus) Subscribe(topic string, handler fnHandler) error {
// 	partitions, err := bus.Subscriber.Partitions(topic)
// 	if err != nil {
// 		return err
// 	}

// 	for _, partition := range partitions {
// 		fmt.Println("PARTITION:", partition)
// 		pc, err := bus.Subscriber.ConsumePartition(topic, partition, sarama.OffsetOldest)
// 		if err != nil {
// 			return err
// 		}

// 		go func() {
// 			for {
// 				select {
// 				case message := <-pc.Messages():
// 					handler(message.Value)
// 				}
// 			}
// 		}()
// 	}

// 	return nil
// }

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/Shopify/sarama"
// )

// type Consumer interface {
// 	Subscribe(handler EventHandler)
// 	Unsubscribe()
// }

// type EventHandler interface {
// 	Setup(sarama.ConsumerGroupSession) error
// 	Cleanup(sarama.ConsumerGroupSession) error
// 	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
// }

// type kafkaConsumer struct {
// 	topics        []string
// 	retryTopic    string
// 	errorTopic    string
// 	consumerGroup sarama.ConsumerGroup
// }

// type ConnectionParameters struct {
// 	Brokers         string
// 	Topics          []string
// 	ErrorTopic      string
// 	RetryTopic      string
// 	ConsumerGroupID string
// }

// func NewConsumer(connectionParams ConnectionParameters) (Consumer, error) {
// 	config := sarama.NewConfig()
// 	cg, err := sarama.NewConsumerGroup(strings.Split(connectionParams.Brokers, ","), connectionParams.ConsumerGroupID, config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &kafkaConsumer{
// 		topics:        connectionParams.Topics,
// 		retryTopic:    connectionParams.RetryTopic,
// 		errorTopic:    connectionParams.ErrorTopic,
// 		consumerGroup: cg,
// 	}, nil
// }

// func (c *kafkaConsumer) Subscribe(handler EventHandler) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	topics := func() []string {
// 		result := make([]string, 0)
// 		if c.errorTopic != "" {
// 			result = append(result, c.errorTopic)
// 		}
// 		if c.retryTopic != "" {
// 			result = append(result, c.retryTopic)
// 		}
// 		result = append(result, c.topics...)
// 		return result
// 	}

// 	go func() {
// 		for {
// 			if err := c.consumerGroup.Consume(ctx, topics(), handler); err != nil {
// 				log.Panicf("Error from consumer : %s", err.Error())
// 			}

// 			if ctx.Err() != nil {
// 				log.Panicf("Error from consumer : %s", ctx.Err().Error())
// 			}
// 		}
// 	}()

// 	go func() {
// 		for err := range c.consumerGroup.Errors() {
// 			fmt.Println("Error from consumer group : ", err.Error())
// 		}
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				fmt.Println("terminating: context cancelled")
// 				cancel()
// 			}
// 		}
// 	}()
// 	fmt.Printf("Kafka consumer listens topics : %v \n", c.topics)
// }

// func (c *kafkaConsumer) Unsubscribe() {
// 	if err := c.consumerGroup.Close(); err != nil {
// 		fmt.Printf("Client wasn't closed :%+v", err)
// 	}
// 	fmt.Println("Kafka consumer closed")
// }
