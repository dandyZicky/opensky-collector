// Package processor contains domain logic related to event processing and consumption from kafka, also includes persistence
package processor

import (
	"context"

	"github.com/dandyZicky/opensky-collector/pkg/events"
	"gorm.io/gorm"
)

type ProcessorService struct {
	Consumer Consumer
	Ctx      context.Context
	DB       *gorm.DB
}

func (p *ProcessorService) NewSubscriberService() {
	p.Consumer.Subscribe(p.Ctx, events.TelemetryRaw)
}

/*
	Websocket vs SSE
	Well, an interface should suffice
*/
