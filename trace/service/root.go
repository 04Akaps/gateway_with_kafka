package service

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	kafkaRepo "github.com/04Akaps/gateway_with_kafka.git/trace/repository/kafka"
	mysqlRepo "github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/handler"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/kafka"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/mysql"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types"
	"log"
	"time"
)

type Service struct {
	kafka   kafka.KafkaImpl
	sql     mysql.MySQLImpl
	handler handler.Handler

	workChan []chan types.WorkChanTrace

	cfg config.TraceCfg
}

func NewService(cfg config.TraceCfg, k *kafkaRepo.Kafka, sql *mysqlRepo.MySQL) *Service {
	s := &Service{
		cfg:      cfg,
		workChan: workChanTraces(len(cfg.Kafka), cfg.ServiceInfo.WorkChanLength),
	}

	var err error

	s.kafka, err = kafka.NewKafkaService(k, s.workChan)

	if err != nil {
		panic(err)
	}

	s.sql = mysql.NewMySQLService(sql)
	s.handler = handler.NewHandler(s.kafka, s.sql)

	return s
}

func (s *Service) PollingAll() {

	index := uint64(0)

	for topic, _ := range s.cfg.Kafka {
		go s.pollingChannel(topic, s.GetWorkChan(index))
		index++
	}

	go s.kafka.PollingEvents()
}

// pollingChannel
// if topic added, need to add new handler in GetHandlerByTopic
func (s *Service) pollingChannel(topic string, workChan chan types.WorkChanTrace) {
	for {
		select {
		case msg, ok := <-workChan:
			if !ok {
				log.Println("Failed to receive channel chan is closed", "time", time.Now().Unix(), "topic", topic)
				continue
			}

			handlerByTopic := s.handler.GetHandlerByTopic(topic)
			handlerByTopic(msg.Response, msg.Message)
		}
	}
}

func (s *Service) GetWorkChan(index uint64) chan types.WorkChanTrace {
	return s.workChan[index]
}

func (s *Service) TotalChannels() int {
	return len(s.workChan)
}

func workChanTraces(totalTopics int, chanLength int64) []chan types.WorkChanTrace {
	channels := make([]chan types.WorkChanTrace, totalTopics)

	for i := 0; i < totalTopics; i++ {
		channels[i] = make(chan types.WorkChanTrace, chanLength)
	}

	return channels
}
