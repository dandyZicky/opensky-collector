package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/domain/processor"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

const (
	ConnTimeoutMs = 5000
	SubTimeoutMs  = 5000
	batchSize     = 100
)

type KafkaConsumer struct {
	Client *kafka.Consumer
	Topic  events.Topic
}

func NewKafkaConsumer(conf *kafka.ConfigMap, topic events.Topic) *KafkaConsumer {
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		log.Panicf("Failed to init kafka consumer client: %s", err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			c.Close()
			panic(r)
		}
	}()

	md, err := c.GetMetadata(nil, true, ConnTimeoutMs)
	if err != nil {
		log.Panicf("Kafka brokers unreachable: %v", err)
	}

	if len(md.Brokers) == 0 {
		log.Panic("No brokers found in cluster metadata")
	}

	log.Printf("Connected to Kafka cluster with %d brokers\n", len(md.Brokers))
	return &KafkaConsumer{
		Client: c,
		Topic:  topic,
	}
}

func (k *KafkaConsumer) Subscribe(ctx context.Context, processor processor.EventProcessor) {
	err := k.Client.Subscribe(k.Topic.String(), nil)
	if err != nil {
		log.Panicf("Subscribing error to kafka topic %s: %s", k.Topic, err.Error())
	}

	run := true
	var batchInputs []events.TelemetryRawEvent
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			ev := k.Client.Poll(SubTimeoutMs)
			switch e := ev.(type) {
			case *kafka.Message:
				event := events.RawMessageToTelemetryRawEvent(e.Value)
				batchInputs = append(batchInputs, event)
			case kafka.Error:
				log.Panicf("Consumer error: %v\n", e)
			default:
				if len(batchInputs) > 0 {
					processor.ProcessEvents(batchInputs, batchSize)
					batchInputs = []events.TelemetryRawEvent{}
				}
			}
		}
	}

	log.Println("Closing consumer...")
	k.Client.Close()
}
