package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/dandyZicky/opensky-collector/dto"
)

type Credentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type FlightClient struct {
	httpClient      *http.Client
	credentials     *Credentials
	authServer      string
	url             string
	accessToken     string
	isAuthenticated bool
	mu              *sync.Mutex
}

func (c *FlightClient) getFlights() (*dto.StatesResponse, error) {
	if !c.isAuthenticated {
		c.mu.Lock()
		defer c.mu.Unlock()
		if err := c.authenticate(); err != nil {
			return nil, err
		}
		c.isAuthenticated = true
	}

	req, err := http.NewRequest("GET", c.url+"/states/all", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.isAuthenticated = false
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.isAuthenticated = false
		return nil, fmt.Errorf("unauthorized access, token may have expired/invalid")
	}

	// var result dto.StatesResponse
	var result map[string]any

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Panic(err)
	}

	states := dto.StatesResponse{}
	states.Time = int64(result["time"].(float64))

	for _, res := range result["states"].([]any) {
		state, err := (*dto.DefaultMapper).ToState(nil, res.([]any))
		if err != nil {
			log.Panic(err)
		}
		states.States = append(states.States, state)
		if len(states.States) >= 5 {
			break
		}
	}

	for _, state := range states.States {
		log.Printf("Flight: %+v\n", *state.Callsign)
	}

	return &states, nil
}

func (c *FlightClient) authenticate() error {
	var err error
	defer func() {
		log.Println("Authentication process completed.")
		if err != nil {
			log.Println("Authentication error:", err)
		}

	}()

	log.Println("Authenticating with credentials:", c.credentials.ClientID)
	data := url.Values{}
	data.Set("client_id", c.credentials.ClientID)
	data.Set("client_secret", c.credentials.ClientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to authenticate, status code:", resp.StatusCode)
		return fmt.Errorf("authentication failed with status: %s", resp.Status)
	}

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Panic(err)
	}

	if token, ok := result["access_token"].(string); ok {
		c.accessToken = token
		c.isAuthenticated = true
		log.Println(("Authenticated successfully, access token obtained."))
		return nil
	}

	err = fmt.Errorf("access_token not found in response")
	return err
}

func poll(ctx context.Context, interval time.Duration, flightClient *FlightClient) {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}
			// Polling logic with retry mechanism
			flights, err := flightClient.getFlights()
			if err != nil {
				for retries := 0; retries < 3; retries++ {
					flights, err = flightClient.getFlights()
					if err == nil {
						break
					}
					time.Sleep(2 * time.Second)
				}
			}
			// TODO: Process flights data -> send to kafka topic
			_ = flights
		}
	}
}

func readCredentials(filePath string) (*Credentials, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var creds Credentials
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

	flightClient := &FlightClient{
		credentials: creds,
		url:         baseURL,
		authServer:  authURL,
		httpClient:  &http.Client{},
		mu:          &sync.Mutex{},
	}

	err = flightClient.authenticate()
	if err != nil {
		log.Panic(err)
	}

	go poll(ctx, tickerInterval, flightClient)
	<-ctx.Done()
}
