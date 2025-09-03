// Package broker contains concrete implementation of event producers
package broker

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/producer"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer(conf *kafka.ConfigMap) *kafka.Producer {
	log.Println("Initializing kafka producer")
	p, err := kafka.NewProducer(conf)
	if err != nil {
		panic("Failed to init kafka producer client")
	}

	md, err := p.GetMetadata(nil, true, 5000)
	if err != nil {
		// No brokers reachable
		log.Fatalf("Kafka brokers unreachable: %v", err)
	}

	if len(md.Brokers) == 0 {
		log.Fatalf("No brokers found in cluster metadata")
	}

	log.Printf("Connected to Kafka cluster with %d brokers\n", len(md.Brokers))
	return p
}

func (k *KafkaProducer) Publish(m producer.Message) error {
	kafkaMessage := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.Topic,
			Partition: kafka.PartitionAny,
		},
		Key:   m.Key(),
		Value: m.Value(),
	}

	err := k.Producer.Produce(&kafkaMessage, nil)
	if err != nil {
		return err
	}
	return nil
}
