package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/AlekSi/pointer"
	"github.com/IBM/sarama"

	conv "telemetry/converter"
	"telemetry/databus/producer"
	"telemetry/repository"
	"telemetry/swagger"
)

const (
	telemetryTopicIn = "telemetry_in"

	commandTelemetryLatest    = "latest"
	commandTelemetryPaginated = "paginated"
)

type TelemetryRepository interface {
	GetLatest(
		ctx context.Context,
		moduleId, delviceId int,
	) (repository.TelemetryValue, error)

	GetPaginated(
		ctx context.Context,
		moduleId, deviceId int,
		limit, offset int,
	) ([]repository.TelemetryValue, bool, error)
}

type Producer interface {
	Produce(ctx context.Context, result producer.CommandTelemetryResult) error
}

type Consumer struct {
	consumer sarama.Consumer
	producer Producer
	repo     TelemetryRepository
}

func New(address string, producer Producer, repo TelemetryRepository) (*Consumer, error) {
	c, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		producer: producer,
		repo:     repo,
	}, nil
}

func (c *Consumer) ProcessTelemetryCommands(ctx context.Context) {
	partConsumer, err := c.consumer.ConsumePartition(telemetryTopicIn, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panicf("failed to create consumer: %v", err)
	}

	defer partConsumer.Close()

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Printf("[consumer] message channel closed")
				return
			}

			var cmd CommandTelemetryIn
			err := json.Unmarshal(msg.Value, &cmd)
			if err != nil {
				log.Printf("[consumer] failed to unmarshal message: %s", err.Error())
				continue
			}

			switch cmd.Action {
			case commandTelemetryLatest:
				value, err := c.repo.GetLatest(ctx, cmd.Device.ModuleId, cmd.Device.Id)
				out := producer.CommandTelemetryResult{
					Action: commandTelemetryLatest,
					Values: []swagger.TelemetryValue{conv.RepoToSwaggerTelemetryValue(value)},
					Err:    conv.ErrorToStringOrNil(err),
				}
				err = c.producer.Produce(ctx, out)
				if err != nil {
					log.Printf("[consumer] failed to produce result for action '%s': %s", cmd.Action, err.Error())
				}
			case commandTelemetryPaginated:
				var perPage int
				var page int
				if cmd.PaginatedParams != nil {
					perPage = pointer.Get(cmd.PaginatedParams.PerPage)
					page = pointer.Get(cmd.PaginatedParams.Page)
				}
				perPage = conv.NormalizePerPage(perPage)

				values, hasMore, err := c.repo.GetPaginated(ctx, cmd.Device.ModuleId, cmd.Device.Id, perPage, page*perPage)
				swaggerOut := conv.RepoToSwaggerPaginatedValues(values, hasMore)
				out := producer.CommandTelemetryResult{
					Action:  commandTelemetryPaginated,
					Values:  swaggerOut.Values,
					HasMore: swaggerOut.HasMore,
					Err:     conv.ErrorToStringOrNil(err),
				}
				err = c.producer.Produce(ctx, out)
				if err != nil {
					log.Printf("[consumer] failed to produce result for action '%s': %s", cmd.Action, err.Error())
				}
			}

		case <-ctx.Done():
			log.Printf("[consumer] context done")
			return
		}
	}
}
