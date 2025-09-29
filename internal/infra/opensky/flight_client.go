package opensky

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dandyZicky/opensky-collector/internal/dto"
	"github.com/dandyZicky/opensky-collector/pkg/retry"
)

type FlightClient struct {
	HTTPClient    *http.Client
	Credentials   *Credentials
	AuthServer    string
	URL           string
	accessToken   string
	authenticated bool
	Mutex         *sync.Mutex
}

func (c *FlightClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *FlightClient) ensureAuthenticated() error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.authenticated {
		return nil
	}

	log.Println("Client not authenticated. Authenticating now...")
	return c.authenticate()
}

func (c *FlightClient) requestAuthorizedStateVectors() (*dto.StatesResponse, error) {
	req, err := http.NewRequest("GET", c.URL+"/states/all", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("lamin", "-11.00")
	q.Set("lomin", "95.00")
	q.Set("lamax", "6.00")
	q.Set("lomax", "141.00")
	req.URL.RawQuery = q.Encode()

	c.Mutex.Lock()
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	c.Mutex.Unlock()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("Token expired or invalid (401 Unauthorized).")
		c.Mutex.Lock()
		c.authenticated = false
		c.Mutex.Unlock()
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	states := dto.StatesResponse{}
	states.Time = int64(result["time"].(float64))
	for _, res := range result["states"].([]any) {
		state, err := (*dto.DefaultMapper).ToState(nil, res.([]any))
		if err != nil {
			log.Printf("Error parsing state vector, skipping: %v", err)
			continue
		}
		states.States = append(states.States, state)
	}

	return &states, nil
}

func (c *FlightClient) GetAllStateVectors() (*dto.StatesResponse, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, fmt.Errorf("initial authentication failed: %w", err)
	}

	states, err := c.requestAuthorizedStateVectors()
	if err != nil {
		if err.Error() == "unauthorized" {
			log.Println("Attempting to re-authenticate and retry the request...")
			authFunc := func() error {
				c.Mutex.Lock()
				defer c.Mutex.Unlock()
				return c.authenticate()
			}

			if authErr := retry.Do(authFunc, retry.WithAttempts(3), retry.WithBackoff(2*time.Second, 2)); authErr != nil {
				return nil, fmt.Errorf("re-authentication failed after multiple attempts: %w", authErr)
			}

			log.Println("Re-authentication successful. Retrying the API request one last time.")
			return c.requestAuthorizedStateVectors()
		}
		return nil, err
	}

	return states, nil
}

func (c *FlightClient) authenticate() error {
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status: %s", resp.Status)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if token, ok := result["access_token"].(string); ok {
		c.accessToken = token
		c.authenticated = true
		log.Println(("Authenticated successfully, access token obtained."))
		return nil
	}

	return fmt.Errorf("access_token not found in response")
}

func ReadCredentials(filePath string) (*Credentials, error) {
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
