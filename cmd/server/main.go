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

	"github.com/dandyZicky/opensky-collector/internal/clients"
	"github.com/dandyZicky/opensky-collector/internal/collector"
)

func readCredentials(filePath string) (*clients.Credentials, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var creds clients.Credentials
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

	flightClient := &clients.FlightClient{
		Credentials: creds,
		URL:         baseURL,
		AuthServer:  authURL,
		HTTPClient:  &http.Client{},
		Mu:          &sync.Mutex{},
	}

	flightDataCollector := &collector.DefaultCollector{
		Client: flightClient,
	}

	err = flightClient.Authenticate()
	if err != nil {
		log.Panic(err)
	}

	go flightDataCollector.Poll(ctx, tickerInterval)
	// go poll(ctx, tickerInterval, flightClient)
	<-ctx.Done()
}
