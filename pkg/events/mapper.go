package events

import (
	"github.com/dandyZicky/opensky-collector/internal/dto"
)

func StateVectorToTelemetryRawEvent(state dto.State) TelemetryRawEvent {
	return TelemetryRawEvent{
		Icao24:        state.Icao24,
		OriginCountry: state.OriginCountry,
		Lat:           nilFloat64(state.Latitude),
		Lon:           nilFloat64(state.Longitude),
		Velocity:      nilFloat64(state.Velocity),
		TimePosition:  nilInt64(state.TimePosition),
		BaroAltitude:  nilFloat64(state.BaroAltitude),
		GeoAltitude:   nilFloat64(state.GeoAltitude),
		LastContact:   state.LastContact,
	}
}

func nilFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func nilInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
