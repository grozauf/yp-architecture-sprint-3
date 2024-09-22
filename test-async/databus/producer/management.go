package producer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

func (p *Producer) ProcessManagementCommands(ctx context.Context) {
	for {
		select {
		case msg := <-p.inputManagementCmd:
			err := p.produceManageCmd(ctx, msg)
			if err != nil {
				log.Printf("[producer] produce manager msg error: %s", err.Error())
			}
		case pError := <-p.producer.Errors():
			log.Printf("[producer] failed to produce manager message: %s", pError.Err.Error())
		case <-ctx.Done():
			log.Printf("[producer] context done")
			return
		}
	}
}

func (p *Producer) ProduceManageCmd(ctx context.Context, cmd ManagementCommand) error {
	select {
	case p.inputManagementCmd <- cmd:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) produceManageCmd(ctx context.Context, cmd ManagementCommand) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: managementTopicIn,
		Value: sarama.ByteEncoder(b),
	}
	select {
	case p.producer.Input() <- &msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
