// Package producer contains logic for producing messages to brokers.
package producer

type Message interface {
	Key() []byte
	Value() []byte
}

type Producer interface {
	Publish(m Message) error
}
