package mqtt_adapter

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/harlesbayu/fleet-backend/internal/domain/model"
	"github.com/harlesbayu/fleet-backend/internal/usecase"
)

type MQTTSubscriber struct {
	client  mqtt.Client
	usecase *usecase.VehicleUsecase
}

func NewMQTTSubscriber(client mqtt.Client, vUsecase *usecase.VehicleUsecase) *MQTTSubscriber {
	return &MQTTSubscriber{
		client:  client,
		usecase: vUsecase,
	}
}

func (s *MQTTSubscriber) SubscribeVehicleLocation() {
	topic := "/fleet/vehicle/+/location"

	handler := func(client mqtt.Client, msg mqtt.Message) {
		var data model.VehicleLocation
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Invalid JSON payload: %v", err)
			return
		}

		if err := data.Validate(); err != nil {
			log.Printf("[MQTT] Invalid vehicle data: %v", err)
			return
		}

		if err := s.usecase.SaveLocation(&data); err != nil {
			log.Printf("Failed to save location: %v", err)
			return
		}

		log.Printf("Saved location for %s (%.4f, %.4f)", data.VehicleID, data.Latitude, data.Longitude)
	}

	if token := s.client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
}
