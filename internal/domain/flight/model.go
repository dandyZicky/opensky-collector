package flight

import (
	"time"

	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type FlightState struct {
	Icao24        string
	OriginCountry string
	Lat           float64
	Lon           float64
	Velocity      float64
	TimePosition  time.Time
	BaroAltitude  float64
	GeoAltitude   float64
	LastContact   time.Time
}

func EventToFlightState(event events.TelemetryRawEvent) FlightState {
	return FlightState{
		Icao24:        event.Icao24,
		OriginCountry: event.OriginCountry,
		Lat:           event.Lat,
		Lon:           event.Lon,
		Velocity:      event.Velocity,
		TimePosition:  time.Unix(event.TimePosition, 0),
		BaroAltitude:  event.BaroAltitude,
		GeoAltitude:   event.GeoAltitude,
		LastContact:   time.Unix(event.LastContact, 0),
	}
}
