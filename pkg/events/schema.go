package events

type TelemetryRawEvent struct {
	Icao24        string  `json:"icao24"`
	OriginCountry string  `json:"origin_country"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Velocity      float64 `json:"velocity"`
	TimePosition  int64   `json:"time_position"`
	BaroAltitude  float64 `json:"baro_altitude"`
	GeoAltitude   float64 `json:"geo_altitude"`
	LastContact   int64   `json:"last_contact"`
}
