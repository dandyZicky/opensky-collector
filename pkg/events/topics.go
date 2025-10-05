package events

type Topic string

var (
	TelemetryRaw      Topic
	TelemetryEnriched Topic
)

func (t Topic) String() string {
	return string(t)
}

func InitTopics(raw, enriched string) {
	TelemetryRaw = Topic(raw)
	TelemetryEnriched = Topic(enriched)
}
