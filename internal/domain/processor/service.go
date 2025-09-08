package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type ProcessorService struct {
	Consumer Consumer
	Ctx      context.Context
}

func (p *ProcessorService) NewSubscriberService() {
	p.Consumer.Subscribe(p.Ctx, events.TelemetryRaw)
}

/*
	Websocket vs SSE
	Well, an interface should suffice
*/
