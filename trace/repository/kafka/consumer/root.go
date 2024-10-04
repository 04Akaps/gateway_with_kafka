package consumer

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

type consumer struct {
	cfg config.KafkaConsumerCfg

	Client *kafka.Consumer

	workChan chan<- types.WorkChanTrace
}

type ConsumerImpl interface {
	PollingEvent()
	CommitMessage(event *kafka.Message) ([]kafka.TopicPartition, error)
}

func NewConsumer(topic string, cfg config.KafkaConsumerCfg, workChan chan types.WorkChanTrace) (ConsumerImpl, error) {
	c := &consumer{
		cfg:      cfg,
		workChan: workChan,
	}

	conf := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.URI,
		"group.id":           cfg.GroupID,
		"auto.offset.reset":  cfg.AutoOffsetReset,
		"enable.auto.commit": false,
	}

	var err error

	c.Client, err = kafka.NewConsumer(conf)

	if err != nil {
		return nil, err
	}

	// we don't need rebalanceCb
	err = c.Client.Subscribe(topic, nil)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *consumer) PollingEvent() {

	defer c.closeHandler()

	for {
		ev := c.Client.Poll(c.cfg.Polling)

		switch event := ev.(type) {
		case *kafka.Message:

			c.workChan <- types.WorkChanTrace{
				Response: event.Value,
				Message:  event,
			}

		case *kafka.Error:
			log.Println("Failed to consume event", "time", time.Now().Unix(), "err", event.Error(), "code", event.Code())
		case *kafka.PartitionEOF:
			log.Println("Failed to consume event partition eof", "time", time.Now().Unix(), "err", event.Error.Error())
		}
	}
}

func (c *consumer) CommitMessage(event *kafka.Message) ([]kafka.TopicPartition, error) {
	return c.Client.CommitMessage(event)
}

func (c *consumer) closeHandler() {
	if err := c.Client.Close(); err != nil {
		log.Println("Failed to close consumer client", "err", err)
	}
}

var _ ConsumerImpl = (*consumer)(nil)
