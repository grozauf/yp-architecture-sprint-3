package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

const (
	telemetryTopicResult = "telemetry_result"
)

type Consumer struct {
	consumer sarama.Consumer
}

func New(address string) (*Consumer, error) {
	c, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
	}, nil
}

func (c *Consumer) ProcessTelemetryCommands(ctx context.Context) {
	partConsumer, err := c.consumer.ConsumePartition(telemetryTopicResult, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panicf("failed to create consumer: %v", err)
	}

	defer partConsumer.Close()

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Printf("[consumer] channel closed, exiting")
				return
			}

			var result CommandTelemetryResult
			err := json.Unmarshal(msg.Value, &result)
			if err != nil {
				log.Printf("[consumer] failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("[consumer] got telemetry result: %v", result)
		case <-ctx.Done():
			log.Printf("[consumer] context done")
			return
		}
	}
}
