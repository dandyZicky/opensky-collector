// Package collector contains domain logic
package collector

import (
	"context"
	"log"
	"time"

	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type CollectorService struct {
	Producer Producer
	Client   Client
}

func (c *CollectorService) Poll(ctx context.Context, interval time.Duration) {
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
			// TODO: Extract processing logic to a function so we can inject it as external dependencies
			flights, err := c.Client.GetAllStateVectors()
			if err != nil {
				unresolved := true
				for range 3 {
					flights, err = c.Client.GetAllStateVectors()
					if err != nil {
						log.Println("Retrying...")
					} else {
						unresolved = false
						break
					}
					time.Sleep(10 * time.Second)
				}
				if unresolved {
					log.Panic("Server couldn't reach OpenSky network")
				}
			}
			// TODO: Process flights data -> send to kafka topic
			for _, state := range flights.States {
				if c.Producer.Publish(events.StateVectorToTelemetryRawEvent(state), events.TelemetryRaw) != nil {
					log.Println("Problems publishing message")
				}
			}
		}
	}
}
