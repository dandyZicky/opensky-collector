// Package collector contains domain logic
package collector

import (
	"context"
	"log"
	"time"

	"github.com/dandyZicky/opensky-collector/pkg/events"
	"github.com/dandyZicky/opensky-collector/pkg/retry"
)

type CollectorService struct {
	Producer Producer
	Client   Client
}

func (c *CollectorService) processAndPublish() error {
	flights, err := c.Client.GetAllStateVectors()
	if err != nil {
		return err
	}

	for _, state := range flights.States {
		if c.Producer.Publish(events.StateVectorToTelemetryRawEvent(state), events.TelemetryRaw) != nil {
			log.Println("Problems publishing message")
		}
	}
	return nil
}

func (c *CollectorService) Poll(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFailures := 0
	const maxConsecutiveFailures = 5 // Escalate after 5 consecutive failed cycles

	for {
		select {
		case <-ctx.Done():
			log.Println("Collector service shutting down.")
			return
		case <-ticker.C:
			err := retry.Do(
				c.processAndPublish,
				retry.WithAttempts(3),
				retry.WithBackoff(10*time.Second, 2.0),
			)

			if err != nil {
				consecutiveFailures++
				log.Printf("Failed cycle %d/%d: %v", consecutiveFailures, maxConsecutiveFailures, err)
				if consecutiveFailures >= maxConsecutiveFailures {
					log.Fatalf("Collector service failed after %d consecutive cycles. Escalating.", maxConsecutiveFailures)
				}
			} else {
				consecutiveFailures = 0
			}
		}
	}
}
