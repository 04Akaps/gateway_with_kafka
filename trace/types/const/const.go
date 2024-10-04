package _const

import _errors "github.com/04Akaps/gateway_with_kafka.git/trace/types/error"

const (
	Report = "report-data"
)

var topics = map[string]bool{
	Report: true,
}

func IsSupportedTopic(topic string) error {
	if _, ok := topics[topic]; !ok {
		return _errors.New(_errors.ErrTopicNotSupported, "topic", topic)
	} else {
		return nil
	}
}
