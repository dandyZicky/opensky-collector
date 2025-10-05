package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/config"
	"github.com/dandyZicky/opensky-collector/internal/domain/collector"
	producer "github.com/dandyZicky/opensky-collector/internal/infra/kafka"
	"github.com/dandyZicky/opensky-collector/internal/infra/opensky"
)

func main() {
	config.InitConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	creds, err := opensky.ReadCredentials(config.AppConfig.OpenSky.CredentialsFile)
	if err != nil {
		log.Fatalf("failed to read credentials: %v", err)
		os.Exit(1)
	}

	flightClient := &opensky.FlightClient{
		Credentials: creds,
		URL:         config.AppConfig.OpenSky.BaseURL,
		AuthServer:  config.AppConfig.OpenSky.AuthURL,
		HTTPClient:  &http.Client{},
		Mutex:       &sync.Mutex{},
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.BootstrapServers,
		"client.id":         config.AppConfig.Kafka.ClientID,
		"acks":              config.AppConfig.Kafka.Acks,
	}

	producerKafka := producer.KafkaProducer{
		Producer: producer.NewKafkaProducer(kafkaConf),
	}

	defer producerKafka.Producer.Close()

	flightDataCollector := &collector.CollectorService{
		Client:   flightClient,
		Producer: &producerKafka,
	}

	go flightDataCollector.Poll(ctx, time.Duration(config.AppConfig.OpenSky.TickerInterval)*time.Millisecond)
	<-ctx.Done()
}
