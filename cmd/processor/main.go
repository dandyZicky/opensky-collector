package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/domain/processor"
	consumer "github.com/dandyZicky/opensky-collector/internal/infra/kafka"
	"github.com/dandyZicky/opensky-collector/internal/infra/pg"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "opensky-processor",
		"auto.offset.reset": "earliest",
	}

	dbConf := pg.NewConfigFromDotEnv(".env")

	db, err := pg.NewDB(dbConf)
	if err != nil {
		log.Panicf("Failed to init db: %s", err.Error())
	}

	kafkaConsumer := consumer.NewKafkaConsumer(kafkaConf, events.TelemetryRaw)
	defer kafkaConsumer.Client.Close()
	flightDataProcessor := &processor.ProcessorService{
		Ctx:      ctx,
		Repo:     &processor.ProcessorRepository{DB: db},
		Consumer: kafkaConsumer,
	}

	go flightDataProcessor.NewSubscriberService()

	<-ctx.Done()

}
