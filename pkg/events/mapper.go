package events

import (
	"encoding/json"

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

func RawMessageToTelemetryRawEvent(rawMessage []byte) TelemetryRawEvent {
	var event TelemetryRawEvent
	err := json.Unmarshal(rawMessage, &event)
	if err != nil {
		return TelemetryRawEvent{}
	}
	return event
}

func SerializeTelemetryRawEvent(event TelemetryRawEvent) ([]byte, error) {
	b, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return b, nil
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
