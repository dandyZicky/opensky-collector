package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/internal/domain/flight"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type EventProcessor interface {
	ProcessEvents(events []events.TelemetryRawEvent, batchSize int) error
}

type Consumer interface {
	Subscribe(ctx context.Context, eventProcessor EventProcessor)
}

type Broadcaster interface {
	Broadcast(events []events.TelemetryRawEvent) error
}

type Inserter interface {
	InsertBatch(states []flight.FlightState, batchSize int) error
}
