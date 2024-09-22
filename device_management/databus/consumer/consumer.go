package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"

	conv "management/converter"
	"management/databus/producer"
	"management/devices"
	"management/swagger"
)

const (
	managementTopicIn = "management_in"

	commandInfo = "info"
	commandSet  = "set"
)

type Producer interface {
	Produce(ctx context.Context, result producer.CommandResult) error
}

type DevicesModuleService interface {
	GetDeviceInfo(ctx context.Context, moduleId, deviceId int) (devices.DeviceInfo, error)
	SetDeviceTargetValue(ctx context.Context, moduleId, deviceId int, value float32) error
	SetDeviceStatus(ctx context.Context, moduleId, deviceId int, status bool) error
}

type Consumer struct {
	consumer  sarama.Consumer
	producer  Producer
	moduleSrv DevicesModuleService
}

func New(address string, producer Producer, moduleSrv DevicesModuleService) (*Consumer, error) {
	c, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:  c,
		producer:  producer,
		moduleSrv: moduleSrv,
	}, nil
}

func (c *Consumer) ProcessCommands(ctx context.Context) {
	partConsumer, err := c.consumer.ConsumePartition(managementTopicIn, 0, sarama.OffsetNewest)
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

			var cmd Command
			err := json.Unmarshal(msg.Value, &cmd)
			if err != nil {
				log.Printf("[consumer] failed to unmarshal message: %s", err.Error())
				continue
			}

			switch cmd.Action {
			case commandInfo:
				info, err := c.moduleSrv.GetDeviceInfo(ctx, cmd.Device.ModuleId, cmd.Device.Id)
				out := producer.CommandResult{
					Action: commandInfo,
					Info: func() *swagger.DeviceInfo {
						info := conv.ToSwaggerDeivceInfo(info)
						return &info
					}(),
					Err: conv.ErrorToStringOrNil(err),
				}
				err = c.producer.Produce(ctx, out)
				if err != nil {
					log.Printf("[consumer] failed to produce result for action '%s': %s", cmd.Action, err.Error())
				}
			case commandSet:
				var err error
				if cmd.Value == nil {
					err = errors.New("empty value not supported in 'set' command")
				} else {
					switch cmd.Value.ValueName {
					case swagger.Status:
						if cmd.Value.Status == nil {
							err = errors.New("empty status not supported in 'set.status' command")
						} else {
							err = c.moduleSrv.SetDeviceStatus(ctx, cmd.Device.ModuleId, cmd.Device.Id, *cmd.Value.Status)
						}
					case swagger.TargetValue:
						if cmd.Value.TargetValue == nil {
							err = errors.New("empty target_value not supported in 'set.target_value' command")
						} else {
							err = c.moduleSrv.SetDeviceTargetValue(ctx, cmd.Device.ModuleId, cmd.Device.Id, *cmd.Value.TargetValue)
						}
					}
				}
				out := producer.CommandResult{
					Action: commandSet,
					Value:  cmd.Value,
					Err:    conv.ErrorToStringOrNil(err),
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
