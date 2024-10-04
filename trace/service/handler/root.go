package handler

import (
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/kafka"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/mysql"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types"
	_const "github.com/04Akaps/gateway_with_kafka.git/trace/types/const"
	"log"
)

type Handler struct {
	kafka kafka.KafkaImpl
	sql   mysql.MySQLImpl

	reportHandler *ReportHandler

	log log.Logger
}

func NewHandler(kafka kafka.KafkaImpl, sql mysql.MySQLImpl) Handler {

	return Handler{
		kafka:         kafka,
		sql:           sql,
		reportHandler: NewReportHandler(kafka, sql),
	}
}

func (s *Handler) GetHandlerByTopic(topic string) types.Handler {

	switch topic {
	case _const.Report:
		return s.reportHandler.handleReport
	default:
		log.Println("Failed to find func by topic when find handler", "topic", topic)
		return nil
	}
}
