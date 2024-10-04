package kafka

import (
	"log"

	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/kafka"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types"

	_error "github.com/04Akaps/gateway_with_kafka.git/trace/types/error"
	_kafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaService struct {
	kafka *kafka.Kafka
}

type KafkaImpl interface {
	PollingEvents()
	CommitMessage(event *_kafka.Message) ([]_kafka.TopicPartition, error)
}

var _ KafkaImpl = (*kafkaService)(nil)

func NewKafkaService(kafka *kafka.Kafka, workChan []chan types.WorkChanTrace) (KafkaImpl, error) {
	s := &kafkaService{
		kafka: kafka,
	}

	err := s.kafka.SetChannelsForConsumeEvent(workChan)

	if err != nil {
		log.Println("Failed to set channel to kafka", "err", err)
		return nil, err
	}

	return s, nil
}

func (k *kafkaService) PollingEvents() {
	for topic := range k.kafka.Consumers {
		k.kafka.Consumers[topic].PollingEvent()
	}
}

func (k *kafkaService) CommitMessage(event *_kafka.Message) ([]_kafka.TopicPartition, error) {
	if consumer, ok := k.kafka.Consumers[*event.TopicPartition.Topic]; !ok {
		return nil, _error.New(_error.ErrConsumerNotFound)
	} else {
		return consumer.CommitMessage(event)
	}
}
