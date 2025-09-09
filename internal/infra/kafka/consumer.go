package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

const (
	ConnTimeoutMs = 5000
	SubTimeoutMs  = 100
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Topic    string
}

func NewKafkaConsumer(conf *kafka.ConfigMap) *kafka.Consumer {
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		log.Panicf("Failed to init kafka consumer client: %s", err.Error())
	}

	md, err := c.GetMetadata(nil, true, ConnTimeoutMs)
	if err != nil {
		log.Fatalf("Kafka brokers unreachable: %v", err)
	}

	if len(md.Brokers) == 0 {
		log.Fatalf("No brokers found in cluster metadata")
	}

	log.Printf("Connected to Kafka cluster with %d brokers\n", len(md.Brokers))
	return c
}

func (k *KafkaConsumer) Subscribe(ctx context.Context, topic events.Topic) {
	err := k.Consumer.Subscribe(string(topic), nil)
	if err != nil {
		log.Panicf("Subscribing error to kafka topic %s: %s", string(topic), err.Error())
	}

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			ev := k.Consumer.Poll(SubTimeoutMs)
			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("Message on %s: %+v\n", e.TopicPartition, events.RawMessageToTelemetryRawEvent(e.Value))
			case kafka.Error:
				log.Panicf("Consumer error: %v\n", e)
			}
		}
	}

	log.Println("Closing consumer...")
	k.Consumer.Close()
}

func (k *KafkaConsumer) insertEventToDB(e *events.TelemetryRawEvent) {
}
