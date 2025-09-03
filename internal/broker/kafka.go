// Package broker contains concrete implementation of event producers
package broker

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/producer"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(conf *kafka.ConfigMap) *KafkaProducer {
	p, err := kafka.NewProducer(conf)
	if err != nil {
		panic("Failed to init kafka producer client")
	}

	return &KafkaProducer{producer: p}
}

func (k *KafkaProducer) Publish(m producer.Message) error {
	kafkaMessage := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   m.Key(),
		Value: m.Value(),
	}

	err := k.producer.Produce(&kafkaMessage, nil)
	if err != nil {
		return err
	}
	return nil
}
