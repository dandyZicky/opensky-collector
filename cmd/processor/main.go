package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/domain/processor"
	consumer "github.com/dandyZicky/opensky-collector/internal/infra/kafka"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "opensky-processor",
		"auto.offset.reset": "earliest",
	}

	consumerKafka := &consumer.KafkaConsumer{
		Consumer: consumer.NewKafkaConsumer(kafkaConf),
	}

	flightDataProcessor := &processor.ProcessorService{
		Consumer: consumerKafka,
		Ctx:      ctx,
	}

	go flightDataProcessor.NewSubscriberService()

}
