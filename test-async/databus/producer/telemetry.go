package producer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

func (p *Producer) ProcessTelemetryCommands(ctx context.Context) {
	for {
		select {
		case msg := <-p.inputTelemetryCmd:
			err := p.produceTelemetryRequest(ctx, msg)
			if err != nil {
				log.Printf("[producer] produce msg error: %s", err.Error())
			}
		case pError := <-p.producer.Errors():
			log.Printf("[producer] failed to produce message: %s", pError.Err.Error())
		case <-ctx.Done():
			log.Printf("[producer] context done")
			p.producer.AsyncClose()
			return
		}
	}
}

func (p *Producer) ProduceTelemetryRequest(ctx context.Context, cmd CommandTelemetryIn) error {
	select {
	case p.inputTelemetryCmd <- cmd:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) produceTelemetryRequest(ctx context.Context, cmd CommandTelemetryIn) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: telemetryTopicIn,
		Value: sarama.ByteEncoder(b),
	}
	select {
	case p.producer.Input() <- &msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
