// Package pg contains timescaledb repository implementation
package pg

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func InsertBatch(db *gorm.DB, flightStateVectors []FlightStateVector, batchSize int) error {
	tx := db.Begin(&sql.TxOptions{})
	err := tx.CreateInBatches(&flightStateVectors, batchSize).Error
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	tx.Commit()
	return nil
}
