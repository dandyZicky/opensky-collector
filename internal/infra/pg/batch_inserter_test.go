package pg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dandyZicky/opensky-collector/internal/domain/flight"
)

func TestPgInserter_InsertBatch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&FlightStateVector{})
	require.NoError(t, err)

	inserter := &PgInserter{DB: db}

	states := []flight.FlightState{
		{
			Icao24:        "test123",
			OriginCountry: "DE",
			Lat:           49.0,
			Lon:           6.0,
			Velocity:      200.0,
			TimePosition:  time.Now(),
			BaroAltitude:  10000,
			GeoAltitude:   10050,
			LastContact:   time.Now(),
		},
		{
			Icao24:        "test456",
			OriginCountry: "FR",
			Lat:           48.0,
			Lon:           2.0,
			Velocity:      250.0,
			TimePosition:  time.Now(),
			BaroAltitude:  12000,
			GeoAltitude:   12050,
			LastContact:   time.Now(),
		},
	}

	err = inserter.InsertBatch(states, 10)

	assert.NoError(t, err)

	var inserted []FlightStateVector
	err = db.Find(&inserted).Error
	assert.NoError(t, err)
	assert.Len(t, inserted, 2)

	assert.Equal(t, "test123", inserted[0].Icao24)
	assert.Equal(t, "DE", inserted[0].OriginCountry)
	assert.Equal(t, 49.0, inserted[0].Lat)
	assert.Equal(t, 6.0, inserted[0].Lon)
	assert.Equal(t, 200.0, inserted[0].Velocity)

	assert.Equal(t, "test456", inserted[1].Icao24)
	assert.Equal(t, "FR", inserted[1].OriginCountry)
	assert.Equal(t, 48.0, inserted[1].Lat)
	assert.Equal(t, 2.0, inserted[1].Lon)
	assert.Equal(t, 250.0, inserted[1].Velocity)
}

func TestPgInserter_InsertBatch_EmptySlice(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&FlightStateVector{})
	require.NoError(t, err)

	inserter := &PgInserter{DB: db}

	err = inserter.InsertBatch([]flight.FlightState{}, 10)

	assert.NoError(t, err)

	var inserted []FlightStateVector
	err = db.Find(&inserted).Error
	assert.NoError(t, err)
	assert.Len(t, inserted, 0)
}

func TestPgInserter_InsertBatch_LargeBatch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&FlightStateVector{})
	require.NoError(t, err)

	inserter := &PgInserter{DB: db}

	states := make([]flight.FlightState, 100)
	for i := range 100 {
		states[i] = flight.FlightState{
			Icao24:        "large" + string(rune(i)),
			OriginCountry: "TEST",
			Lat:           float64(i),
			Lon:           float64(i * 2),
			Velocity:      float64(i * 10),
			TimePosition:  time.Now(),
			BaroAltitude:  float64(i * 100),
			GeoAltitude:   float64(i*100 + 50),
			LastContact:   time.Now(),
		}
	}

	err = inserter.InsertBatch(states, 25)

	assert.NoError(t, err)

	var count int64
	err = db.Model(&FlightStateVector{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(100), count)
}

func TestPgInserter_InsertBatch_DatabaseError(t *testing.T) {
	inserter := &PgInserter{DB: nil} // if DB was nil test

	states := []flight.FlightState{
		{
			Icao24:        "error123",
			OriginCountry: "DE",
			Lat:           49.0,
			Lon:           6.0,
			Velocity:      200.0,
			TimePosition:  time.Now(),
			BaroAltitude:  10000,
			GeoAltitude:   10050,
			LastContact:   time.Now(),
		},
	}

	err := inserter.InsertBatch(states, 10)

	assert.Error(t, err)
}
