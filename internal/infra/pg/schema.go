package pg

import (
	"time"
)

type FlightStateVector struct {
	ID            uint      `gorm:"primaryKey"`
	OriginCountry string    `gorm:"not null"`
	Lat           float64   `gorm:"not null"`
	Lon           float64   `gorm:"not null"`
	Velocity      float64   `gorm:"not null"`
	TimePosition  time.Time `gorm:"type:timestamp not null"`
	BaroAltitude  float64   `gorm:"not null"`
	GeoAltitude   float64   `gorm:"not null"`
	LastContact   time.Time `gorm:"type:timestamp not null"`
}
