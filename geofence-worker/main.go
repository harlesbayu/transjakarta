package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

// Config
const (
	rabbitMQUrl = "amqp://guest:guest@rabbitmq:5672/"
)

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Timestamp int64 `json:"timestamp"`
}

func main() {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange & queue (sama seperti di fleet-backend)
	err = ch.ExchangeDeclare(
		"fleet.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare(
		"geofence_alerts",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,
		"geofence_alerts",
		"fleet.events",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to consume messages: %v", err)
	}

	// Handle messages asynchronously
	log.Println("[Worker] waiting for geofence events...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for msg := range msgs {
			var event GeofenceEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("[Worker] failed to unmarshal event: %v", err)
				continue
			}

			log.Printf("[Worker] Vehicle %s triggered %s at (%.6f, %.6f) time=%d",
				event.VehicleID,
				event.Event,
				event.Location.Latitude,
				event.Location.Longitude,
				event.Timestamp,
			)
		}
	}()

	<-sigCh
	log.Println("[Worker] shutting down gracefully...")
}
