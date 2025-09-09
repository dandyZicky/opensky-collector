package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type Consumer interface {
	Subscribe(ctx context.Context, topic events.Topic)
}

type Repository interface {
	Save(ctx context.Context, events []events.TelemetryRawEvent)
}
