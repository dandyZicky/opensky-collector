// Package collector contains implementations of data polling
package collector

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/dandyZicky/opensky-collector/internal/broker"
	"github.com/dandyZicky/opensky-collector/internal/clients"
	"github.com/dandyZicky/opensky-collector/internal/dto"
	"github.com/dandyZicky/opensky-collector/internal/producer"
	"github.com/dandyZicky/opensky-collector/pkg/events"
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
			flights, err := c.Client.GetAllStateVectors()
			if err != nil {
				for range 3 {
					flights, err = c.Client.GetAllStateVectors()
					if err == nil {
						log.Println("Retrying...")
						break
					}
					time.Sleep(2 * time.Second)
				}
			}
			// TODO: Process flights data -> send to kafka topic
			for _, state := range flights.States {
				msg, err := c.toMessage(state)
				if err != nil {
					log.Println("Failed to convert to kafka message")
					continue
				}
				if c.Producer.Publish(msg) != nil {
					log.Println("Problems publishing message")
				}
			}
		}
	}
}

func (c *DefaultCollector) toMessage(state dto.State) (producer.Message, error) {
	telemetryEvent := events.StateVectorToTelemetryRawEvent(state)
	payload, err := json.Marshal(telemetryEvent)
	if err != nil {
		return producer.DefaultMessage{}, err
	}

	msg := producer.DefaultMessage{
		MessageKey:   telemetryEvent.Icao24,
		MessageValue: payload,
	}

	return msg, nil
}
