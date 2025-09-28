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
	"github.com/dandyZicky/opensky-collector/internal/domain/collector"
	producer "github.com/dandyZicky/opensky-collector/internal/infra/kafka"
	"github.com/dandyZicky/opensky-collector/internal/infra/opensky"
)

const (
	baseURL        = "https://opensky-network.org/api"
	authURL        = "https://auth.opensky-network.org/auth/realms/opensky-network/protocol/openid-connect/token"
	tickerInterval = 21600 * time.Millisecond
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	creds, err := opensky.ReadCredentials("credentials.json")
	if err != nil {
		log.Fatalf("failed to read credentials: %v", err)
		os.Exit(1)
	}

	flightClient := &opensky.FlightClient{
		Credentials: creds,
		URL:         baseURL,
		AuthServer:  authURL,
		HTTPClient:  &http.Client{},
		Mu:          &sync.Mutex{},
	}

	err = flightClient.Authenticate()
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
		os.Exit(1)
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"client.id":         "openskyCollector",
		"acks":              "all",
	}

	producerKafka := producer.KafkaProducer{
		Producer: producer.NewKafkaProducer(kafkaConf),
	}

	defer producerKafka.Producer.Close()

	flightDataCollector := &collector.CollectorService{
		Client:   flightClient,
		Producer: &producerKafka,
	}

	go flightDataCollector.Poll(ctx, tickerInterval)
	// go poll(ctx, tickerInterval, flightClient)
	<-ctx.Done()
}
