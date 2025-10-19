package helper

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
)

// GeofenceLocation represents a single geofence area
type GeofenceLocation struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

// cache the JSON load result to avoid re-reading the file
var (
	geofenceLocations []GeofenceLocation
	geofenceOnce      sync.Once
)

// LoadGeofenceLocations reads the geofence_locations.json
func LoadGeofenceLocations() []GeofenceLocation {
	geofenceOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			log.Printf("[Geofence] failed to get working directory: %v", err)
			return
		}

		filePath := filepath.Join(wd, "configs", "geofence_locations.json")

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[Geofence] failed to read file %s: %v", filePath, err)
			return
		}

		if err := json.Unmarshal(data, &geofenceLocations); err != nil {
			log.Printf("[Geofence] failed to unmarshal JSON: %v", err)
			return
		}

		log.Printf("[Geofence] loaded %d geofence locations from %s", len(geofenceLocations), filePath)
	})
	return geofenceLocations
}

// CheckGeofence checks whether the point (lat, lon) is within the radius of any location
// and returns (bool, location name)
func CheckGeofence(lat, lon float64) (bool, string) {
	locations := LoadGeofenceLocations()

	for _, loc := range locations {
		distance := haversineDistance(lat, lon, loc.Latitude, loc.Longitude)
		if distance <= loc.Radius {
			log.Printf("[Geofence] inside area: %s (distance %.2fm)", loc.Name, distance)
			return true, loc.Name
		}
	}
	return false, ""
}

// haversineDistance calculates the distance between two points in meters
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0 // meter

	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)
	lat1Rad := toRadians(lat1)
	lat2Rad := toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
