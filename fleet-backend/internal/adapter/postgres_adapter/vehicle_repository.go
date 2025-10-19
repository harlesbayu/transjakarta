package postgres_adapter

import (
	"github.com/harlesbayu/fleet-backend/internal/domain/model"
	"github.com/harlesbayu/fleet-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) repository.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Save(loc *model.VehicleLocation) error {
	return r.db.Create(loc).Error
}

func (r *vehicleRepository) GetLastLocation(vehicleID string) (*model.VehicleLocation, error) {
	var loc model.VehicleLocation
	err := r.db.Where("vehicle_id = ?", vehicleID).Order("timestamp DESC").First(&loc).Error
	return &loc, err
}

func (r *vehicleRepository) GetHistory(vehicleID string, start, end int64) ([]model.VehicleLocation, error) {
	var list []model.VehicleLocation
	err := r.db.Where("vehicle_id = ? AND timestamp BETWEEN ? AND ?", vehicleID, start, end).Order("timestamp asc").Find(&list).Error
	return list, err
}
