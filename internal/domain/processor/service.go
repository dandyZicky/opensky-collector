// Package processor contains domain logic related to event processing and consumption from kafka, also includes persistence
package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/internal/infra/pg"
	"github.com/dandyZicky/opensky-collector/pkg/events"
	"gorm.io/gorm"
)

type ProcessorService struct {
	Consumer    Consumer
	Ctx         context.Context
	Repo        FlightStateRepository
	Broadcaster Broadcaster
}

type ProcessorRepository struct {
	DB *gorm.DB
}

func (p *ProcessorService) NewSubscriberService() {
	processor := p
	p.Consumer.Subscribe(p.Ctx, processor)
}

func (p *ProcessorService) ProcessEvents(events []events.TelemetryRawEvent, batchSize int) error {
	var states []pg.FlightStateVector
	p.Broadcaster.Broadcast(events)
	for _, event := range events {
		states = append(states, pg.EventToFlightStateVector(event))
	}
	return p.Repo.SaveBatch(states, batchSize)
}

func (r *ProcessorRepository) SaveBatch(states []pg.FlightStateVector, batchSize int) error {
	return pg.InsertBatch(r.DB, states, batchSize)
}
