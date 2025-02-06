package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer

func initKafka(bootstrapServers string) {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
}

func sendToKafka(count int64) {
	topic := "unique_request_counts"
	message := fmt.Sprintf("Unique request count: %d", count)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	}
}
