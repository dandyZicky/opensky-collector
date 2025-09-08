// Package kafka contains concrete implementation of event producers
package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/pkg/events"
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

func (k *KafkaProducer) Publish(event events.TelemetryRawEvent, topic events.Topic) error {
	msg, err := EventToMessage(event, string(topic))
	if err != nil {
		return err

	}
	err = k.Producer.Produce(msg, nil)
	if err != nil {
		return err
	}
	return nil
}
