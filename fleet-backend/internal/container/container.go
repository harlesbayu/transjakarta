package container

import (
	"log"

	mqttAdapter "github.com/harlesbayu/fleet-backend/internal/adapter/mqtt_adapter"
	postgresAdapter "github.com/harlesbayu/fleet-backend/internal/adapter/postgres_adapter"
	rabbitmqAdapter "github.com/harlesbayu/fleet-backend/internal/adapter/rabbitmq_adapter"
	"github.com/harlesbayu/fleet-backend/internal/config"
	"github.com/harlesbayu/fleet-backend/internal/handler"
	"github.com/harlesbayu/fleet-backend/internal/usecase"
)

type Container struct {
	DB             *postgresAdapter.DBWrapper
	Config         *config.Config
	VehicleHandler *handler.VehicleHandler
	PingHandler    *handler.PingHandler
}

func NewContainer(db *postgresAdapter.DBWrapper, cfg *config.Config) *Container {
	// Init RabbitMQ
	rmq, err := rabbitmqAdapter.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("failed to connect RabbitMQ: %v", err)
	}

	vRepo := postgresAdapter.NewVehicleRepository(db.DB)
	vUsecase := usecase.NewVehicleUsecase(vRepo, rmq)
	vHandler := handler.NewVehicleHandler(vUsecase)

	mqttClient := mqttAdapter.NewMQTTClient(cfg.MQTT.Host, cfg.MQTT.Port)
	mqttSub := mqttAdapter.NewMQTTSubscriber(mqttClient, vUsecase)
	mqttSub.SubscribeVehicleLocation()

	return &Container{
		DB:             db,
		Config:         cfg,
		VehicleHandler: vHandler,
		PingHandler:    handler.NewPingHandler(),
	}
}
