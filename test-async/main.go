package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"

	"test_async/databus/consumer"
	prdcr "test_async/databus/producer"
)

const (
	KAFKA_ADDRESS_ENV = "KAFKA"
	KAFKA_ADDRESS     = "kafka:9092"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault(KAFKA_ADDRESS_ENV, KAFKA_ADDRESS)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaAddress := viper.GetString(KAFKA_ADDRESS_ENV)

	producer, err := prdcr.New(kafkaAddress)
	if err != nil {
		log.Panicf("failed to start kafka producer: %s\n", err.Error())
	}
	go producer.ProcessTelemetryResults(ctx)

	consumer, err := consumer.New(kafkaAddress)
	if err != nil {
		log.Panicf("failed to start kafka consumer: %s\n", err.Error())
	}
	go consumer.ProcessTelemetryCommands(ctx)

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		log.Printf("[test_async] produce telemetry cmd '%s'\n", prdcr.CommandTelemetryLatest)
		err = producer.Produce(ctx, prdcr.CommandTelemetryIn{
			Action: prdcr.CommandTelemetryLatest,
			Device: prdcr.Device{
				Id:       0,
				ModuleId: 0,
			},
		})
		if err != nil {
			log.Printf("produce cmd failed: %s\n", err)
		}

		select {
		case <-signals:
			log.Printf("[test_async] finished")
			return
		default:
		}

		time.Sleep(time.Second * 2)
	}
}
