package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"telemetry/handler"
	"telemetry/repository"
	"telemetry/swagger"
)

const (
	SERVER_PORT_ENV   = "PORT"
	SERVER_PORT       = "8080"
	KAFKA_ADDRESS_ENV = "KAFKA"
	KAFKA_ADDRESS     = "kafka:9092"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault(SERVER_PORT_ENV, SERVER_PORT)
	viper.SetDefault(KAFKA_ADDRESS_ENV, KAFKA_ADDRESS)

	repo := repository.New()

	handler := handler.New(repo)

	r := gin.Default()
	swagger.RegisterHandlers(r, handler)
	r.Run(":" + viper.GetString(SERVER_PORT_ENV))
}
