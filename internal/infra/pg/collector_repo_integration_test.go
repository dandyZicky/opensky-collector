// Package pg_test contains timescaledb repository integration tests
package pg

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"

	_ "embed"

	"github.com/dandyZicky/opensky-collector/pkg/events"
	"github.com/joho/godotenv"
)

//go:embed testdata/.env.db.test
var envTest []byte

const batchSize = 1000

func loadConfigFromTestEnv() Config {
	envMap, err := godotenv.UnmarshalBytes(envTest)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return Config{
		Host:     envMap["POSTGRES_HOST"],
		User:     envMap["POSTGRES_USER"],
		Password: envMap["POSTGRES_PASSWORD"],
		Port:     envMap["POSTGRES_PORT"],
		Dbname:   envMap["POSTGRES_DB"],
	}
}

func loadRawTelemetryData(filepath string) []events.TelemetryRawEvent {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var events []events.TelemetryRawEvent
	err = json.NewDecoder(file).Decode(&events)
	if err != nil {
		panic(err)
	}

	return events
}

func TestInsertBatchFlightStateVector(t *testing.T) {
	db, err := NewDB(loadConfigFromTestEnv())
	if err != nil {
		t.Fatal(err)
	}

	events := loadRawTelemetryData("test_data.json")

	var msgs []FlightStateVector
	for _, v := range events {
		msgs = append(msgs, EventToFlightStateVector(v))
	}

	err = insertBatch(db, msgs, batchSize)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success")

}

func BenchmarkInsertBatchFlightStateVector(b *testing.B) {
	db, err := NewDB(loadConfigFromTestEnv())
	if err != nil {
		b.Fatal(err)
	}

	events := loadRawTelemetryData("test_data.json")

	var msgs []FlightStateVector
	for _, v := range events {
		msgs = append(msgs, EventToFlightStateVector(v))
	}

	for b.Loop() {
		err = insertBatch(db, msgs, batchSize)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInsertFlightStateVector(b *testing.B) {
	db, err := NewDB(loadConfigFromTestEnv())
	if err != nil {
		b.Fatal(err)
	}

	events := loadRawTelemetryData("test_data.json")

	var msgs []FlightStateVector
	for _, v := range events {
		msgs = append(msgs, EventToFlightStateVector(v))
	}

	for b.Loop() {
		tx := db.Begin(&sql.TxOptions{})
		for _, v := range msgs {
			tx.Create(&v)
			if err != nil {
				b.Fatal(err)
			}
		}
		tx.Rollback()
	}
}
