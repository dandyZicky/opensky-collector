package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/internal/infra/pg"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type EventProcessor interface {
	ProcessEvents(events []events.TelemetryRawEvent, batchSize int) error
}

type FlightStateRepository interface {
	SaveBatch(states []pg.FlightStateVector, batchSize int) error
}

type Consumer interface {
	Subscribe(ctx context.Context, eventProcessor EventProcessor)
}

type Broadcaster interface {
	Broadcast(events []events.TelemetryRawEvent) error
}
