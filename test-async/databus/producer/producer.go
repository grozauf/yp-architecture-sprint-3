package producer

import (
	"github.com/IBM/sarama"
)

const (
	telemetryTopicIn  = "telemetry_in"
	managementTopicIn = "management_in"
)

type Producer struct {
	producer           sarama.AsyncProducer
	inputTelemetryCmd  chan CommandTelemetryIn
	inputManagementCmd chan ManagementCommand
}

func New(address string) (*Producer, error) {
	p, err := sarama.NewAsyncProducer([]string{address}, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer:           p,
		inputTelemetryCmd:  make(chan CommandTelemetryIn, 10),
		inputManagementCmd: make(chan ManagementCommand, 10),
	}, nil
}
