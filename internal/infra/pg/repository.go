package pg

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InsertFlightStateVector(flightStateVector *FlightStateVector) error {
	return nil
}

func NewDB() (*gorm.DB, error) {
	config := newConfig()
	log.Printf("DSN: %+v", config)
	dsn := "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.Dbname + " port=" + config.Port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
