package usecase

import (
	"context"
	"log"
	"time"

	rabbitmqAdapter "github.com/harlesbayu/fleet-backend/internal/adapter/rabbitmq_adapter"
	"github.com/harlesbayu/fleet-backend/internal/domain/model"
	"github.com/harlesbayu/fleet-backend/internal/domain/repository"
	"github.com/harlesbayu/fleet-backend/internal/helper"
)

type VehicleUsecase struct {
	repo repository.VehicleRepository
	rmq  *rabbitmqAdapter.RabbitMQ
}

func NewVehicleUsecase(r repository.VehicleRepository, rmq *rabbitmqAdapter.RabbitMQ) *VehicleUsecase {
	return &VehicleUsecase{repo: r, rmq: rmq}
}

func (u *VehicleUsecase) SaveLocation(loc *model.VehicleLocation) error {
	if err := u.repo.Save(loc); err != nil {
		return err
	}

	u.handleGeofenceEvent(loc)
	return nil
}

func (u *VehicleUsecase) GetLastLocation(vehicleID string) (*model.VehicleLocation, error) {
	return u.repo.GetLastLocation(vehicleID)
}

func (u *VehicleUsecase) GetHistory(vehicleID string, from, to int64) ([]model.VehicleLocation, error) {
	return u.repo.GetHistory(vehicleID, from, to)
}

func (u *VehicleUsecase) handleGeofenceEvent(v *model.VehicleLocation) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Geofence] panic recovered: %v", r)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		inside, locationName := helper.CheckGeofence(v.Latitude, v.Longitude)

		if inside {
			log.Printf("[Geofence] ✅ Vehicle %s is INSIDE %s (lat=%.6f, lon=%.6f)",
				v.VehicleID, locationName, v.Latitude, v.Longitude)
		} else {
			log.Printf("[Geofence] ❌ Vehicle %s is OUTSIDE geofence (lat=%.6f, lon=%.6f)",
				v.VehicleID, v.Latitude, v.Longitude)
			return
		}

		done := make(chan error, 1)
		go func() {
			done <- u.rmq.PublishGeofenceEvent(v.VehicleID, v.Latitude, v.Longitude, v.Timestamp)
		}()

		select {
		case err := <-done:
			if err != nil {
				log.Printf("[RabbitMQ] publish error: %v", err)
			} else {
				log.Printf("[RabbitMQ] event published for %s (%s)", v.VehicleID, locationName)
			}
		case <-ctx.Done():
			log.Printf("[RabbitMQ] publish timeout for %s", v.VehicleID)
		}
	}()
}
