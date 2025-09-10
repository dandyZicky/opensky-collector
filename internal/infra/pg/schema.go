package pg

import (
	"time"

	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type FlightStateVector struct {
	ID            uint      `gorm:"primaryKey"`
	Icao24        string    `gorm:"not null"`
	OriginCountry string    `gorm:"not null"`
	Lat           float64   `gorm:"not null"`
	Lon           float64   `gorm:"not null"`
	Velocity      float64   `gorm:"not null"`
	TimePosition  time.Time `gorm:"type:timestamp not null"`
	BaroAltitude  float64   `gorm:"not null"`
	GeoAltitude   float64   `gorm:"not null"`
	LastContact   time.Time `gorm:"type:timestamp not null"`
}

func EventToFlightStateVector(event events.TelemetryRawEvent) FlightStateVector {
	return FlightStateVector{
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
