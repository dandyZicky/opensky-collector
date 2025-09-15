// Package pg contains timescaledb repository implementation
package pg

import (
	"database/sql"
	"fmt"

	"github.com/dandyZicky/opensky-collector/internal/domain/flight"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgInserter struct {
	DB *gorm.DB
}

func NewDB(config Config) (*gorm.DB, error) {
	dsn := "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.Dbname + " port=" + config.Port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&FlightStateVector{}); err != nil {
		return nil, err
	}
	return db, nil
}

func (p *PgInserter) InsertBatch(flightState []flight.FlightState, batchSize int) error {
	if p.DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	tx := p.DB.Begin(&sql.TxOptions{})

	var flightStateVectors []FlightStateVector
	for _, state := range flightState {
		flightStateVectors = append(flightStateVectors, ToFlightStateVector(state))
	}

	err := tx.CreateInBatches(&flightStateVectors, batchSize).Error
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	tx.Commit()
	return nil
}
