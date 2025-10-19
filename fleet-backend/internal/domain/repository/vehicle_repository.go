package repository

import "github.com/harlesbayu/fleet-backend/internal/domain/model"

type VehicleRepository interface {
	Save(loc *model.VehicleLocation) error
	GetLastLocation(vehicleID string) (*model.VehicleLocation, error)
	GetHistory(vehicleID string, start, end int64) ([]model.VehicleLocation, error)
}
