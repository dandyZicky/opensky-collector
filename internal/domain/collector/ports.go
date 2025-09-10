package collector

import (
	"context"
	"net/http"
	"time"

	"github.com/dandyZicky/opensky-collector/internal/dto"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type Message interface {
	Key() []byte
	Value() []byte
}

type Producer interface {
	Publish(event events.TelemetryRawEvent, topic events.Topic) error
}

type Collector interface {
	Poll(ctx context.Context, interval time.Duration)
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	GetAllStateVectors() (*dto.StatesResponse, error)
}
