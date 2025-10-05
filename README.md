# OpenSky Flight Monitor

## Project Overview

The OpenSky Flight Monitor is a backend system designed to collect, process, persist, and broadcast real-time flight telemetry data from the OpenSky Network API. It's built with a microservice architecture, featuring a dedicated data collector and a real-time data processor, communicating via Apache Kafka. Data is persisted in a PostgreSQL database optimized for time-series data using the TimescaleDB extension, and real-time updates are pushed to connected clients via Server-Sent Events (SSE). This project is intended to be paired with a frontend application for visualization and interaction.

## Features

*   **Real-time Data Collection:** Continuously polls the OpenSky Network API for flight state vectors.
*   **Event-Driven Architecture:** Utilizes Apache Kafka for asynchronous and decoupled communication between services.
*   **Data Processing & Persistence:** Consumes raw flight data, transforms it, and stores it efficiently in PostgreSQL with TimescaleDB.
*   **Real-time Broadcast:** Pushes live flight updates to connected frontend clients using Server-Sent Events (SSE).
*   **Robust API Integration:** Handles OpenSky Network API authentication (OAuth2 client credentials) and includes retry mechanisms for resilience.
*   **Centralized Configuration:** Manages application settings via `config.yaml` and environment variables using `spf13/viper`.
*   **Graceful Shutdown:** Ensures clean termination of services upon interruption signals.

## Architecture

The system is composed of two main Go services:

1.  **`collector` Service:**
    *   Responsible for making authenticated requests to the OpenSky Network API.
    *   Fetches all available flight state vectors within a predefined geographical bounding box.
    *   Publishes these raw flight telemetry events to a Kafka topic (`telemetry.raw`).

2.  **`processor` Service:**
    *   Subscribes to the `telemetry.raw` Kafka topic.
    *   Consumes raw flight events, transforms them into a more suitable domain model.
    *   Persists the processed flight state data into a PostgreSQL database (with TimescaleDB).
    *   Simultaneously broadcasts the real-time flight updates to connected clients via an SSE endpoint.

**Data Flow:**
`OpenSky API` &rarr; `Collector Service` &rarr; `Kafka (telemetry.raw)` &rarr; `Processor Service` &rarr; `PostgreSQL/TimescaleDB`
`Processor Service` &rarr; `SSE Broadcaster` &rarr; `Frontend Clients`

## Technologies Used

*   **Backend Language:** Go
*   **Message Broker:** Apache Kafka (`github.com/confluentinc/confluent-kafka-go/v2`)
*   **Database:** PostgreSQL with TimescaleDB extension
*   **External API:** OpenSky Network API
*   **Real-time Communication:** Server-Sent Events (SSE)
*   **Configuration Management:** `spf13/viper`
*   **Concurrency:** Go routines, `context.Context`, `sync.Mutex`
*   **HTTP Client:** Standard Go `net/http`

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   Go (1.22+)
*   Docker & Docker Compose (for Kafka and PostgreSQL)
*   An OpenSky Network API account (for client ID and secret)

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/dandyZicky/opensky-collector.git
    cd opensky-collector
    ```

2.  **Configure Environment:**
    *   Create a `credentials.json` file in the root directory (or specified by `opensky.credentials_file` in `config.yaml`) with your OpenSky API client ID and secret:
        ```json
        {
          "client_id": "YOUR_OPENSKY_CLIENT_ID",
          "client_secret": "YOUR_OPENSKY_CLIENT_SECRET"
        }
        ```
    *   Review and modify `internal/config/config.yaml` as needed. This file contains default configurations for the database, Kafka, SSE, and OpenSky API. You can override these settings using environment variables (e.g., `KAFKA_BOOTSTRAP_SERVERS=localhost:9092`).

3.  **Start Infrastructure Services (Kafka, PostgreSQL):**
    Use Docker Compose to spin up the required infrastructure:
    ```bash
    docker-compose up -d
    ```
    This will start Kafka, Zookeeper, and PostgreSQL with TimescaleDB.

### Running the Application

1.  **Run the Collector Service:**
    ```bash
    go run cmd/collector/main.go
    ```

2.  **Run the Processor Service:**
    ```bash
    go run cmd/processor/main.go
    ```

Both services will start, and the collector will begin fetching data, publishing it to Kafka. The processor will consume this data, store it, and broadcast it via SSE.

## Frontend Integration

The `processor` service exposes an SSE endpoint for real-time flight data. Your frontend application can connect to this endpoint to receive live updates. The default endpoint is `http://localhost:8081/sse/flights`. Ensure your frontend's origin is listed in `sse.allowed_origins` in `config.yaml`.