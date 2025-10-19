package mqtt_adapter

import (
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	client mqtt.Client
	once   sync.Once
)

// NewMQTTClient initializes the MQTT client (singleton)
func NewMQTTClient(host string, port int) mqtt.Client {
	once.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
		opts.SetClientID("fleet-backend-subscriber")
		opts.SetAutoReconnect(true)
		opts.OnConnect = func(c mqtt.Client) {
			log.Println("Connected to MQTT broker")
		}
		opts.OnConnectionLost = func(c mqtt.Client, err error) {
			log.Printf("Connection lost: %v", err)
		}

		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
		}
	})

	return client
}
