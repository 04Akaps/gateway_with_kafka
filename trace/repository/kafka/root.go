package kafka

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/kafka/consumer"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types"
	_const "github.com/04Akaps/gateway_with_kafka.git/trace/types/const"
)

type Kafka struct {
	Consumers map[string]consumer.ConsumerImpl

	cfg config.TraceCfg
}

func NewKafka(cfg config.TraceCfg) (*Kafka, error) {
	k := &Kafka{
		cfg:       cfg,
		Consumers: make(map[string]consumer.ConsumerImpl),
	}

	return k, nil
}

func (k *Kafka) SetChannelsForConsumeEvent(workChan []chan types.WorkChanTrace) error {
	// we don't need atomic
	index := 0
	var err error

	for topic, kafkaConf := range k.cfg.Kafka {
		if err = _const.IsSupportedTopic(topic); err != nil {
			return err
		}

		k.Consumers[topic], err = consumer.NewConsumer(topic, kafkaConf, workChan[index])

		if err != nil {
			return err
		}

		index++
	}

	return nil
}
