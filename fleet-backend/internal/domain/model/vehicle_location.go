package model

import (
	"errors"
	"fmt"
	"regexp"
)

// VehicleLocation represents a vehicle's current position.
type VehicleLocation struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func (v *VehicleLocation) Validate() error {
	if v.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}

	validVehicleID := regexp.MustCompile(`^[A-Z]{1,2}\d{1,4}[A-Z]{1,3}$`)
	if !validVehicleID.MatchString(v.VehicleID) {
		return fmt.Errorf("invalid vehicle_id format: %s", v.VehicleID)
	}

	if v.Latitude < -90 || v.Latitude > 90 {
		return fmt.Errorf("invalid latitude: %.6f (must be between -90 and 90)", v.Latitude)
	}

	if v.Longitude < -180 || v.Longitude > 180 {
		return fmt.Errorf("invalid longitude: %.6f (must be between -180 and 180)", v.Longitude)
	}

	if !v.isNearJakarta() {
		return fmt.Errorf("coordinates (%.6f, %.6f) appear outside Jakarta area", v.Latitude, v.Longitude)
	}

	return nil
}

func (v *VehicleLocation) isNearJakarta() bool {
	const (
		minLat = -6.4
		maxLat = -6.0
		minLon = 106.6
		maxLon = 107.1
	)
	return v.Latitude >= minLat && v.Latitude <= maxLat &&
		v.Longitude >= minLon && v.Longitude <= maxLon
}
