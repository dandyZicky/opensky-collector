package events

type Topic string

const (
	TelemetryRaw      Topic = "telemetry.raw"
	TelemetryEnriched Topic = "telemetry.enriched"
)

func (t Topic) String() string {
	return string(t)
}
