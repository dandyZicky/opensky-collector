package processor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dandyZicky/opensky-collector/internal/domain/flight"
	"github.com/dandyZicky/opensky-collector/pkg/events"
)

type MockInserter struct {
	mock.Mock
}

func (m *MockInserter) InsertBatch(states []flight.FlightState, batchSize int) error {
	args := m.Called(states, batchSize)
	return args.Error(0)
}

type MockBroadcaster struct {
	mock.Mock
}

func (m *MockBroadcaster) Broadcast(events []events.TelemetryRawEvent) error {
	args := m.Called(events)
	return args.Error(0)
}

type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) Subscribe(ctx context.Context, eventProcessor EventProcessor) {
	m.Called(ctx, eventProcessor)
}

func TestProcessorService_ProcessEvents_Success(t *testing.T) {
	mockInserter := &MockInserter{}
	mockBroadcaster := &MockBroadcaster{}

	processor := &ProcessorService{
		Inserter:    mockInserter,
		Broadcaster: mockBroadcaster,
	}

	events := []events.TelemetryRawEvent{
		{
			Icao24:        "abc123",
			OriginCountry: "DE",
			Lat:           49.0,
			Lon:           6.0,
			Velocity:      200.0,
			TimePosition:  1638360000,
			BaroAltitude:  10000,
			GeoAltitude:   10050,
			LastContact:   1638360000,
		},
	}
	batchSize := 10

	mockBroadcaster.On("Broadcast", events).Return(nil)
	mockInserter.On("InsertBatch", mock.MatchedBy(func(states []flight.FlightState) bool {
		return len(states) == 1 &&
			states[0].Icao24 == "abc123" &&
			states[0].OriginCountry == "DE"
	}), batchSize).Return(nil)

	err := processor.ProcessEvents(events, batchSize)

	assert.NoError(t, err)
	mockInserter.AssertExpectations(t)
	mockBroadcaster.AssertExpectations(t)
}

func TestProcessorService_ProcessEvents_InserterError(t *testing.T) {
	mockInserter := &MockInserter{}
	mockBroadcaster := &MockBroadcaster{}

	processor := &ProcessorService{
		Inserter:    mockInserter,
		Broadcaster: mockBroadcaster,
	}

	events := []events.TelemetryRawEvent{
		{
			Icao24:        "def456",
			OriginCountry: "FR",
			Lat:           48.0,
			Lon:           2.0,
			Velocity:      250.0,
			TimePosition:  1638360000,
			BaroAltitude:  12000,
			GeoAltitude:   12050,
			LastContact:   1638360000,
		},
	}
	batchSize := 5
	expectedError := errors.New("database connection failed")

	mockBroadcaster.On("Broadcast", events).Return(nil)
	mockInserter.On("InsertBatch", mock.AnythingOfType("[]flight.FlightState"), batchSize).Return(expectedError)

	err := processor.ProcessEvents(events, batchSize)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockInserter.AssertExpectations(t)
	mockBroadcaster.AssertExpectations(t)
}

func TestProcessorService_ProcessEvents_BroadcasterError(t *testing.T) {
	mockInserter := &MockInserter{}
	mockBroadcaster := &MockBroadcaster{}

	processor := &ProcessorService{
		Inserter:    mockInserter,
		Broadcaster: mockBroadcaster,
	}

	events := []events.TelemetryRawEvent{
		{
			Icao24:        "ghi789",
			OriginCountry: "GB",
			Lat:           51.0,
			Lon:           0.0,
			Velocity:      300.0,
			TimePosition:  1638360000,
			BaroAltitude:  15000,
			GeoAltitude:   15050,
			LastContact:   1638360000,
		},
	}
	batchSize := 8
	expectedError := errors.New("SSE broadcast failed")

	mockBroadcaster.On("Broadcast", events).Return(expectedError)

	err := processor.ProcessEvents(events, batchSize)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockBroadcaster.AssertExpectations(t)
	mockInserter.AssertNotCalled(t, "InsertBatch", mock.Anything, mock.Anything)
}

func TestProcessorService_ProcessEvents_EmptyEvents(t *testing.T) {
	mockInserter := &MockInserter{}
	mockBroadcaster := &MockBroadcaster{}

	processor := &ProcessorService{
		Inserter:    mockInserter,
		Broadcaster: mockBroadcaster,
	}

	events := []events.TelemetryRawEvent{}
	batchSize := 10

	mockBroadcaster.On("Broadcast", events).Return(nil)
	mockInserter.On("InsertBatch", mock.MatchedBy(func(states []flight.FlightState) bool {
		return len(states) == 0
	}), batchSize).Return(nil)

	err := processor.ProcessEvents(events, batchSize)

	assert.NoError(t, err)
	mockInserter.AssertExpectations(t)
	mockBroadcaster.AssertExpectations(t)
}

func TestProcessorService_ProcessEvents_MultipleEvents(t *testing.T) {
	mockInserter := &MockInserter{}
	mockBroadcaster := &MockBroadcaster{}

	processor := &ProcessorService{
		Inserter:    mockInserter,
		Broadcaster: mockBroadcaster,
	}

	events := []events.TelemetryRawEvent{
		{
			Icao24:        "multi1",
			OriginCountry: "US",
			Lat:           40.0,
			Lon:           -74.0,
			Velocity:      400.0,
			TimePosition:  1638360000,
			BaroAltitude:  20000,
			GeoAltitude:   20050,
			LastContact:   1638360000,
		},
		{
			Icao24:        "multi2",
			OriginCountry: "CA",
			Lat:           45.0,
			Lon:           -75.0,
			Velocity:      350.0,
			TimePosition:  1638360000,
			BaroAltitude:  18000,
			GeoAltitude:   18050,
			LastContact:   1638360000,
		},
	}
	batchSize := 15

	mockBroadcaster.On("Broadcast", events).Return(nil)
	mockInserter.On("InsertBatch", mock.MatchedBy(func(states []flight.FlightState) bool {
		return len(states) == 2 &&
			states[0].Icao24 == "multi1" &&
			states[1].Icao24 == "multi2"
	}), batchSize).Return(nil)

	err := processor.ProcessEvents(events, batchSize)

	assert.NoError(t, err)
	mockInserter.AssertExpectations(t)
	mockBroadcaster.AssertExpectations(t)
}

func TestProcessorService_NewSubscriberService(t *testing.T) {
	mockInserter := &MockInserter{}
	mockBroadcaster := &MockBroadcaster{}
	mockConsumer := &MockConsumer{}

	ctx := context.Background()

	processor := &ProcessorService{
		Inserter:    mockInserter,
		Broadcaster: mockBroadcaster,
		Consumer:    mockConsumer,
		Ctx:         ctx,
	}

	mockConsumer.On("Subscribe", ctx, processor).Return()

	processor.NewSubscriberService()

	mockConsumer.AssertExpectations(t)
}
