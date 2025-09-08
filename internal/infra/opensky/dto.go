package opensky

type StatesResponse struct {
	Time   int64   `json:"time"`
	States []State `json:"states"`
}

type State struct {
	Icao24         string   `json:"icao24"`
	Callsign       *string  `json:"callsign"`
	OriginCountry  string   `json:"origin_country"`
	TimePosition   *int64   `json:"time_position"`
	LastContact    int64    `json:"last_contact"`
	Longitude      *float64 `json:"longitude"`
	Latitude       *float64 `json:"latitude"`
	BaroAltitude   *float64 `json:"baro_altitude"`
	OnGround       bool     `json:"on_ground"`
	Velocity       *float64 `json:"velocity"`
	TrueTrack      *float64 `json:"true_track"`
	VerticalRate   *float64 `json:"vertical_rate"`
	Sensors        []int    `json:"sensors"`
	GeoAltitude    *float64 `json:"geo_altitude"`
	Squawk         *string  `json:"squawk"`
	Spi            bool     `json:"spi"`
	PositionSource int      `json:"position_source"`
	Category       int      `json:"category"`
}

type DefaultMapper struct{}

func (m *DefaultMapper) ToState(data []any) (State, error) {
	if len(data) < 17 {
		return State{}, nil // or return an error indicating insufficient data
	}

	state := State{}

	if v, ok := data[0].(string); ok {
		state.Icao24 = v
	}

	if v, ok := data[1].(string); ok {
		state.Callsign = &v
	}

	if v, ok := data[2].(string); ok {
		state.OriginCountry = v
	}

	if v, ok := data[3].(float64); ok {
		t := int64(v)
		state.TimePosition = &t
	}

	if v, ok := data[4].(float64); ok {
		state.LastContact = int64(v)
	}

	if v, ok := data[5].(float64); ok {
		state.Longitude = &v
	}

	if v, ok := data[6].(float64); ok {
		state.Latitude = &v
	}

	if v, ok := data[7].(float64); ok {
		state.BaroAltitude = &v
	}

	if v, ok := data[8].(bool); ok {
		state.OnGround = v
	}

	if v, ok := data[9].(float64); ok {
		state.Velocity = &v
	}

	if v, ok := data[10].(float64); ok {
		state.TrueTrack = &v
	}

	if v, ok := data[11].(float64); ok {
		state.VerticalRate = &v
	}

	if v, ok := data[12].([]any); ok {
		sensors := make([]int, len(v))
		for i, s := range v {
			if sensorID, ok := s.(float64); ok {
				sensors[i] = int(sensorID)
			}
		}
		state.Sensors = sensors
	}

	if v, ok := data[13].(float64); ok {
		state.GeoAltitude = &v
	}

	if v, ok := data[14].(string); ok {
		state.Squawk = &v
	}

	if v, ok := data[15].(bool); ok {
		state.Spi = v
	}

	if v, ok := data[16].(float64); ok {
		state.PositionSource = int(v)
	}

	if len(data) > 17 {
		if v, ok := data[17].(float64); ok {
			state.Category = int(v)
		}
	}

	return state, nil
}
