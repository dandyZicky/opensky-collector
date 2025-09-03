// Package collector contains implementations of data polling
package collector

import (
	"context"
	"log"
	"time"

	"github.com/dandyZicky/opensky-collector/internal/broker"
	"github.com/dandyZicky/opensky-collector/internal/clients"
	"github.com/dandyZicky/opensky-collector/internal/producer"
)

type Collector interface {
	Poll(ctx context.Context, interval time.Duration, client clients.Client)
}

type DefaultCollector struct {
	Client   *clients.FlightClient
	Producer broker.KafkaProducer
}

func (c *DefaultCollector) Poll(ctx context.Context, interval time.Duration) {
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
			flights, err := c.Client.GetFlights()
			if err != nil {
				for range 3 {
					flights, err = c.Client.GetFlights()
					if err == nil {
						log.Println("Retrying...")
						break
					}
					time.Sleep(2 * time.Second)
				}
			}
			// TODO: Process flights data -> send to kafka topic
			_ = flights
			msg := producer.DefaultMessage{
				MessageKey:   "mantap",
				MessageValue: "jaya",
			}

			if c.Producer.Publish(msg) != nil {
				log.Println("Problems publishing message")
			}
		}
	}
}
