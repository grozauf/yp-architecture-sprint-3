package producer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

const (
	managementTopicResult = "management_result"
)

type Producer struct {
	producer sarama.AsyncProducer
	input    chan CommandResult
}

func New(address string) (*Producer, error) {
	p, err := sarama.NewAsyncProducer([]string{address}, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
		input:    make(chan CommandResult, 10),
	}, nil
}

func (p *Producer) ProcessResults(ctx context.Context) {
	for {
		select {
		case msg := <-p.input:
			err := p.produce(ctx, msg)
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

func (p *Producer) Produce(ctx context.Context, result CommandResult) error {
	select {
	case p.input <- result:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) produce(ctx context.Context, result CommandResult) error {
	b, err := json.Marshal(result)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: managementTopicResult,
		Value: sarama.ByteEncoder(b),
	}
	select {
	case p.producer.Input() <- &msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
