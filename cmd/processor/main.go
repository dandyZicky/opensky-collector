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
	"github.com/dandyZicky/opensky-collector/internal/infra/sse"
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

	inserter := pg.PgInserter{DB: db}

	broadcasterSSE := sse.NewSSEBroadcaster(ctx)
	sseServer := sse.NewSSEServer(broadcasterSSE, "8081")
	go broadcasterSSE.Run()
	go sseServer.Start()

	kafkaConsumer := consumer.NewKafkaConsumer(kafkaConf, events.TelemetryRaw)
	defer kafkaConsumer.Client.Close()
	flightDataProcessor := &processor.ProcessorService{
		Ctx:         ctx,
		Inserter:    &inserter,
		Consumer:    kafkaConsumer,
		Broadcaster: broadcasterSSE,
	}

	go flightDataProcessor.NewSubscriberService()

	<-ctx.Done()

}
