package types

import (
	"github.com/04Akaps/gateway_with_kafka.git/trace/types/request"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types/response"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Handler func(msg []byte, kafkaMsg *kafka.Message)

// TODO -> Publish하는 타입에 맞춰서 작성 및 수정 가능
type TraceEvent struct {
	ServiceID string  `json:"service_id"`
	Latency   float64 `json:"latency"`

	Request  request.Request   `json:"request"`
	Response response.Response `json:"response"`

	StartTime int64 `json:"start_time"`
}

type WorkChanTrace struct {
	Response []byte
	Message  *kafka.Message
}
