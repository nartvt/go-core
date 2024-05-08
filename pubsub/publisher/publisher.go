package publisher

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type NATSPublisher struct {
	nc    *nats.Conn
	topic string
}

func NewPublisher(host string, topic string) *NATSPublisher {
	// Connect to the NATS server
	//nc, err := nats.Connect("nats://localhost:4222")
	nc, err := nats.Connect(host)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// Create a publisher
	publisher := NewNATSPublisher(nc, topic)

	return publisher
}

// NATSPublisher is a simple wrapper around the NATS connection for publishing messages.

// NewNATSPublisher creates a new NATSPublisher.
func NewNATSPublisher(nc *nats.Conn, topic string) *NATSPublisher {
	return &NATSPublisher{nc: nc, topic: topic}
}

// Publish publishes a message to a given subject.
func (p *NATSPublisher) Publish(subject string, msg interface{}) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.nc.Publish(subject, msgByte)
}

// Close closes the NATS connection (optional).
func (p *NATSPublisher) Close() {
	p.nc.Close()
}
