package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Config
const (
	mqttBroker        = "tcp://mosquitto:1883"
	publishInterval   = 2 * time.Second
	insideProbability = 0.6 // 1.0 = 100% pasti di dalam geofence
	safeRadiusFactor  = 0.8 // gunakan max 80% dari radius agar aman dari pembulatan
)

// Geofence location
var geofenceLocations = []struct {
	Name      string
	Latitude  float64
	Longitude float64
	Radius    float64
}{
	{"Monas", -6.175392, 106.827153, 50},
	{"Dukuh Atas 2", -6.201224, 106.822977, 50},
	{"Harmoni Central", -6.166104, 106.829361, 50},
	{"Blok M", -6.244273, 106.800706, 50},
	{"Kampung Melayu", -6.222821, 106.868919, 50},
	{"Grogol 1", -6.157930, 106.788350, 50},
	{"Ragunan", -6.303207, 106.820502, 50},
	{"Pulo Gadung 2", -6.189308, 106.901060, 50},
	{"Kalideres", -6.159167, 106.706389, 50},
	{"Depo Cijantung", -6.317290, 106.852281, 50},
}

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

// func randomVehicleID(r *rand.Rand) string {
// 	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
// 	num := r.Intn(9000) + 1000 // 1000–9999
// 	return fmt.Sprintf("B%d%c%c", num,
// 		letters[r.Intn(len(letters))],
// 		letters[r.Intn(len(letters))],
// 	)
// }

func randomVehicleID(r *rand.Rand) string {
	vehicleIDs := []string{
		"B1001AA",
		"B1002AB",
		"B1003AC",
		"B1004AD",
		"B1005AE",
		"B1006AF",
		"B1007AG",
		"B1008AH",
		"B1009AI",
		"B1010AJ",
	}
	return vehicleIDs[r.Intn(len(vehicleIDs))]
}

func randomNear(r *rand.Rand, lat, lng, radiusMeters float64) (float64, float64) {
	radiusMeters *= safeRadiusFactor
	rnd := math.Sqrt(r.Float64())
	theta := r.Float64() * 2 * math.Pi
	dLat := (radiusMeters * rnd * math.Cos(theta)) / 111_320.0
	dLng := (radiusMeters * rnd * math.Sin(theta)) / (111_320.0 * math.Cos(lat*math.Pi/180))
	return lat + dLat, lng + dLng
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	clientID := fmt.Sprintf("mock-publisher-%d", time.Now().UnixNano())
	opts := mqtt.NewClientOptions().
		AddBroker(mqttBroker).
		SetClientID(clientID).
		SetCleanSession(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect MQTT: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Printf("[MockPublisher] Connected to %s", mqttBroker)
	log.Printf("[MockPublisher] Inside probability: %.0f%%", insideProbability*100)

	for {
		var lat, lng float64
		vehicleID := randomVehicleID(r)

		if r.Float64() < insideProbability {
			spot := geofenceLocations[r.Intn(len(geofenceLocations))]
			lat, lng = randomNear(r, spot.Latitude, spot.Longitude, spot.Radius)
		} else {
			lat = -6.2 + (r.Float64()-0.5)*0.2
			lng = 106.8 + (r.Float64()-0.5)*0.2
		}

		data := VehicleLocation{
			VehicleID: vehicleID,
			Latitude:  lat,
			Longitude: lng,
			Timestamp: time.Now().Unix(),
		}

		payload, _ := json.Marshal(data)
		topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)

		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		log.Printf("[MockPublisher] Sent %s: %.6f, %.6f", vehicleID, lat, lng)
		time.Sleep(publishInterval)
	}
}
