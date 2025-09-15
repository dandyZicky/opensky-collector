package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dandyZicky/opensky-collector/pkg/events"
	"github.com/rs/cors"
)

type SSEBroadcaster struct {
	clients    map[chan []events.TelemetryRawEvent]bool
	register   chan chan []events.TelemetryRawEvent
	unregister chan chan []events.TelemetryRawEvent
	messages   chan []events.TelemetryRawEvent
	ctx        context.Context
}

type SSEServer struct {
	broadcaster *SSEBroadcaster
	port        string
}

func NewSSEServer(broadcaster *SSEBroadcaster, port string) *SSEServer {
	return &SSEServer{
		broadcaster: broadcaster,
		port:        port,
	}
}

func (s *SSEServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/sse/flights", s.handleSSE)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Printf("SSE Server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, handler)
}

func (s *SSEServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	ch := s.broadcaster.Join()
	s.broadcaster.ServeSSE(w, r, ch)
}

func NewSSEBroadcaster(ctx context.Context) *SSEBroadcaster {
	return &SSEBroadcaster{
		clients:    make(map[chan []events.TelemetryRawEvent]bool),
		register:   make(chan chan []events.TelemetryRawEvent, 10),
		unregister: make(chan chan []events.TelemetryRawEvent, 10),
		messages:   make(chan []events.TelemetryRawEvent, 100),
		ctx:        ctx,
	}
}

func (b *SSEBroadcaster) Run() {
	for {
		select {
		case <-b.ctx.Done():
			for ch := range b.clients {
				log.Println("Closing active clients...")
				close(ch)
			}
			return
		case ch := <-b.register:
			b.clients[ch] = true
			log.Println("Registered new session")
		case ch := <-b.unregister:
			log.Println("A session left")
			delete(b.clients, ch)
			close(ch)
		case msgs := <-b.messages:
			for client := range b.clients {
				select {
				case client <- msgs:
				default:
					log.Println("Dropped")
				}
			}
		}
	}
}

func (b *SSEBroadcaster) Join() chan []events.TelemetryRawEvent {
	messageChannel := make(chan []events.TelemetryRawEvent)
	b.register <- messageChannel
	return messageChannel
}

func (b *SSEBroadcaster) Leave(client chan []events.TelemetryRawEvent) {
	b.unregister <- client
}

func (b *SSEBroadcaster) Broadcast(event []events.TelemetryRawEvent) error {
	b.messages <- event
	return nil
}

func (b *SSEBroadcaster) ServeSSE(w http.ResponseWriter, r *http.Request, ch chan []events.TelemetryRawEvent) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Client goroutine panic: %v", r)
			b.Leave(ch)
		}
	}()

	w.Header().Set("content-type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	confirmation := map[string]interface{}{
		"status":    "connected",
		"type":      "connection_ack",
		"timestamp": time.Now().Unix(),
		"message":   "Flight tracking stream active",
	}
	msg, _ := json.Marshal(confirmation)
	fmt.Fprintf(w, "data: %s\n\n", msg)
	w.(http.Flusher).Flush()

	go func() {
		<-r.Context().Done()
		b.Leave(ch)
	}()

	for evs := range ch {
		for _, event := range evs {
			msg, err := events.SerializeTelemetryRawEvent(event)
			if err != nil {
				log.Printf("Failed to serialize event: %v", err)
			}

			fmt.Fprintf(w, "data: %s\n\n", string(msg))
			w.(http.Flusher).Flush()
		}
	}
}
