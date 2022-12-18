package messaging

import (
	"fmt"
	"log"
	"testing"
)

func logHandler1(payload interface{}) error {
	fmt.Println("logHandler1: payload:", payload)

	return nil
}

func logHandler2(payload interface{}) error {
	fmt.Println("logHandler2: payload:", payload)

	return nil
}

func logHandler3(payload interface{}) error {
	fmt.Println("logHandler3: payload:", payload)

	return nil
}

func logHandler4(payload interface{}) error {
	fmt.Println("logHandler4: payload:", payload)

	return nil
}

func logHandler5(payload interface{}) error {
	fmt.Println("logHandler5: payload:", payload)

	return nil
}

// It should publish round robin to partitions of the topic and listen rounde robin
func TestOneTopicMultipleSubscribersSameGroup(t *testing.T) {
	var url = "localhost:9092"
	var groupId = "testing"
	var topic = "my-topic"
	var publishRetries = 3

	pub := NewPublisher(url, publishRetries)
	defer pub.Close()

	for i := 0; i < 5; i++ {
		message := Message{
			Payload: i,
		}

		err := pub.Publish(topic, message)

		if err != nil {
			log.Panicln(err)
		}
	}

	sub := NewSubscriber(url, groupId)

	go sub.Subscribe(topic, logHandler1)
	go sub.Subscribe(topic, logHandler2)
	go sub.Subscribe(topic, logHandler3)
	go sub.Subscribe(topic, logHandler4)
	go sub.Subscribe(topic, logHandler5)

	for {
	}
}

// All consumers should listen to all events
// func TestOneTopicMultipleSubscribersDifferentGroup(t *testing.T) {
// 	var url = "localhost:9092"
// 	var groupId_1 = "testing1"
// 	var groupId_2 = "testing2"
// 	var groupId_3 = "testing3"
// 	var groupId_4 = "testing4"
// 	var groupId_5 = "testing5"
// 	var topic = "my-topic"
// 	var publishRetries = 3

// 	pub := NewPublisher(url, publishRetries)
// 	defer pub.Close()

// 	for i := 0; i < 5; i++ {
// 		message := Message{
// 			Payload: i,
// 		}

// 		err := pub.Publish(topic, message)

// 		if err != nil {
// 			log.Panicln(err)
// 		}
// 	}

// 	sub_1 := NewSubscriber(url, groupId_1)
// 	go sub_1.Subscribe(topic, logHandler1)

// 	sub_2 := NewSubscriber(url, groupId_2)
// 	go sub_2.Subscribe(topic, logHandler2)

// 	sub_3 := NewSubscriber(url, groupId_3)
// 	go sub_3.Subscribe(topic, logHandler3)

// 	sub_4 := NewSubscriber(url, groupId_4)
// 	go sub_4.Subscribe(topic, logHandler4)

// 	sub_5 := NewSubscriber(url, groupId_5)
// 	go sub_5.Subscribe(topic, logHandler5)

// 	// for {
// 	// }
// }
