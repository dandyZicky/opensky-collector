// Package pg contains timescaledb repository implementation
package pg

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InsertFlightStateVector(flightStateVector *FlightStateVector) error {
	return nil
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

func insertBatch(db *gorm.DB, flightStateVectors []FlightStateVector, batchSize int) error {
	tx := db.CreateInBatches(&flightStateVectors, batchSize)
	if err := tx.Error; err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}
