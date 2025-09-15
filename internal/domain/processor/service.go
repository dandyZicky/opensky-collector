// Package processor contains domain logic related to event processing and consumption from kafka, also includes persistence
package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/internal/domain/flight"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type ProcessorService struct {
	Inserter    Inserter
	Consumer    Consumer
	Ctx         context.Context
	Broadcaster Broadcaster
}

func (p *ProcessorService) NewSubscriberService() {
	processor := p
	p.Consumer.Subscribe(p.Ctx, processor)
}

func (p *ProcessorService) ProcessEvents(events []events.TelemetryRawEvent, batchSize int) error {
	var states []flight.FlightState

	// Broadcast events first
	if err := p.Broadcaster.Broadcast(events); err != nil {
		return err
	}

	// Convert events to domain models
	for _, event := range events {
		states = append(states, flight.EventToFlightState(event))
	}

	// Insert flight states
	return p.Inserter.InsertBatch(states, batchSize)
}
