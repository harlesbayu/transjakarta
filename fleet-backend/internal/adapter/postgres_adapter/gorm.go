package postgres_adapter

import (
	"fmt"
	"sync"

	"github.com/harlesbayu/fleet-backend/internal/config"
	"github.com/harlesbayu/fleet-backend/internal/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBWrapper struct {
	DB *gorm.DB
}

var (
	dbInstance *DBWrapper
	once       sync.Once
	dbErr      error
)

// NewGormDB (Singleton)
func NewGormDB(cfg config.PostgresConfig) (*DBWrapper, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			dbErr = err
			return
		}

		// Auto migrate table
		if err := db.AutoMigrate(&model.VehicleLocation{}); err != nil {
			dbErr = err
			return
		}

		dbInstance = &DBWrapper{DB: db}
	})

	return dbInstance, dbErr
}
