package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/dandyZicky/opensky-collector/internal/dto"
)

type FlightClient struct {
	HTTPClient      *http.Client
	Credentials     *Credentials
	AuthServer      string
	URL             string
	AccessToken     string
	IsAuthenticated bool
	Mu              *sync.Mutex
}

func (c *FlightClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *FlightClient) GetAllStateVectors() (*dto.StatesResponse, error) {
	if !c.IsAuthenticated {
		c.Mu.Lock()
		defer c.Mu.Unlock()
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
		c.IsAuthenticated = true
	}

	req, err := http.NewRequest("GET", c.URL+"/states/all", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("lamin", "45.8389")
	q.Set("lomin", "5.9962")
	q.Set("lamax", "47.8229")
	q.Set("lomax", "10.5226")

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.Do(req)
	if err != nil {
		c.IsAuthenticated = false
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.IsAuthenticated = false
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
		// if len(states.States) >= 5 {
		// 	break
		// }
	}

	// for _, state := range states.States {
	// 	log.Printf("Flight: %+v\n", *state.Callsign)
	// }

	return &states, nil
}

func (c *FlightClient) Authenticate() error {
	var err error
	defer func() {
		log.Println("Authentication process completed.")
		if err != nil {
			log.Println("Authentication error:", err)
		}

	}()

	log.Println("Authenticating with credentials:", c.Credentials.ClientID)
	data := url.Values{}
	data.Set("client_id", c.Credentials.ClientID)
	data.Set("client_secret", c.Credentials.ClientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", c.AuthServer, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
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
		c.AccessToken = token
		c.IsAuthenticated = true
		log.Println(("Authenticated successfully, access token obtained."))
		return nil
	}

	err = fmt.Errorf("access_token not found in response")
	return err
}
