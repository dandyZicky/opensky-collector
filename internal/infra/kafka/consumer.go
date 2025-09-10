package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/infra/pg"
	"github.com/dandyZicky/opensky-collector/pkg/events"
	"gorm.io/gorm"
)

const (
	ConnTimeoutMs = 5000
	SubTimeoutMs  = 5000
	batchSize     = 100
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Topic    string
	DB       *gorm.DB
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
	var batchInputs []pg.FlightStateVector
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			ev := k.Consumer.Poll(SubTimeoutMs)
			switch e := ev.(type) {
			case *kafka.Message:
				rawEvent := events.RawMessageToTelemetryRawEvent(e.Value)
				batchInputs = append(batchInputs, pg.EventToFlightStateVector(rawEvent))
			case kafka.Error:
				log.Panicf("Consumer error: %v\n", e)
			default:
				if len(batchInputs) > 0 {
					pg.InsertBatch(k.DB, batchInputs, batchSize)
					batchInputs = []pg.FlightStateVector{}
				}
			}
		}
	}

	log.Println("Closing consumer...")
	k.Consumer.Close()
}
