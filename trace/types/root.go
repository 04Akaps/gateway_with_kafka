package types

import "github.com/confluentinc/confluent-kafka-go/kafka"

type WorkChanTrace struct {
	Response []byte
	Message  *kafka.Message
}
