package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"telemetry/databus/consumer"
	"telemetry/databus/producer"
	"telemetry/handler"
	"telemetry/repository"
	"telemetry/swagger"
)

const (
	SERVER_PORT_ENV = "PORT"
	SERVER_PORT     = "8080"

	KAFKA_ENABLED_ENV = "KAFKA_ENABLED"
	KAFKA_ENABLED     = "true"
	KAFKA_ADDRESS_ENV = "KAFKA"
	KAFKA_ADDRESS     = "kafka:9092"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault(SERVER_PORT_ENV, SERVER_PORT)
	viper.SetDefault(KAFKA_ADDRESS_ENV, KAFKA_ADDRESS)
	viper.SetDefault(KAFKA_ENABLED_ENV, KAFKA_ENABLED)

	repo := repository.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	isKafkaEnabled := viper.GetBool(KAFKA_ENABLED_ENV)
	kafkaAddress := viper.GetString(KAFKA_ADDRESS_ENV)

	log.Printf("message broker address: %s, enable: %v\n", kafkaAddress, isKafkaEnabled)

	if isKafkaEnabled {
		producer, err := producer.New(kafkaAddress)
		if err != nil {
			log.Panicf("failed to start kafka producer: %s\n", err.Error())
		}
		go producer.ProcessTelemetryResults(ctx)

		consumer, err := consumer.New(kafkaAddress, producer, repo)
		if err != nil {
			log.Panicf("failed to start kafka consumer: %s\n", err.Error())
		}
		go consumer.ProcessTelemetryCommands(ctx)
	}

	router := gin.Default()
	swagger.RegisterHandlers(router, handler.New(repo))
	srv := &http.Server{
		Addr:    ":" + viper.GetString(SERVER_PORT_ENV),
		Handler: router.Handler(),
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		if err := srv.Shutdown(ctx); err != nil {
			log.Panicf("gracefull shutdown failed: %s\n", err.Error())
		}
	}()

	log.Printf("starting server on '%s'\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panicf("server listen failed: %s\n", err)
	}

	log.Println("server stopped")
}
