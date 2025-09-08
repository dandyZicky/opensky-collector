package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dandyZicky/opensky-collector/internal/domain/collector"
	broker "github.com/dandyZicky/opensky-collector/internal/infra/kafka"
	"github.com/dandyZicky/opensky-collector/internal/infra/opensky"
)

func readCredentials(filePath string) (*opensky.Credentials, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var creds opensky.Credentials
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&creds)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}

const (
	baseURL        = "https://opensky-network.org/api"
	authURL        = "https://auth.opensky-network.org/auth/realms/opensky-network/protocol/openid-connect/token"
	tickerInterval = 10 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	creds, err := readCredentials("credentials.json")
	if err != nil {
		panic(err)
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
		log.Panic(err)
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"client.id":         "openskyCollector",
		"acks":              "all",
	}

	producerKafka := broker.KafkaProducer{
		Producer: broker.NewKafkaProducer(kafkaConf),
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
