package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

func EventToMessage(e events.TelemetryRawEvent, topic string) (*kafka.Message, error) {
	val, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(e.Icao24),
		Value: val,
	}, nil
}
