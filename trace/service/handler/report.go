package handler

import (
	"encoding/json"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/kafka"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service/mysql"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types"
	_report "github.com/04Akaps/gateway_with_kafka.git/trace/types/report"
	_kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/upper/db/v4"
	"log"
	"sync"
	"time"
)

const (
	_HOUR = "HOUR"
)

type ReportHandler struct {
	kafka kafka.KafkaImpl
	sql   mysql.MySQLImpl

	processLock sync.Mutex
}

func NewReportHandler(kafka kafka.KafkaImpl, sql mysql.MySQLImpl) *ReportHandler {
	handler := &ReportHandler{
		kafka: kafka,
		sql:   sql,
	}

	return handler
}

func (s *ReportHandler) handleReport(value []byte, msg *_kafka.Message) {
	s.processLock.Lock()
	defer s.processLock.Unlock()

	var events []*types.TraceEvent

	if err := json.Unmarshal(value, &events); err != nil {
		log.Println("Failed to unmarshal event", "topic", *msg.TopicPartition.Topic, "value", value, "err", err)
		return
	}

	eventMapper := make(map[string]*_report.Report)

	for _, event := range events {
		ID := event.Request.ID

		if v, ok := eventMapper[ID]; !ok {
			eventMapper[ID] = getReportWhenIDNotFound(event)
		} else {
			eventMapper[ID] = getReportWhenIDFounded(v, event)
		}

	}

	err := s.sql.TxContext(func(tx db.Session) error {

		if err := s.updateDatabase(tx, eventMapper); err != nil {
			return err
		}

		if err := s.commitOffsets(msg); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println("Failed to handle report event", "err", err)
	}
}

func (s *ReportHandler) updateDatabase(tx db.Session, mapper map[string]*_report.Report) error {
	dbResult, err := s.sql.UpdateReportUsingBulkWithTx(tx, mapper)

	if err != nil {
		return err
	}

	affected, _ := dbResult.RowsAffected()

	log.Println("Success update report in handler", "totalMappers", len(mapper), "affected", affected)

	return nil
}

func (s *ReportHandler) commitOffsets(msg *_kafka.Message) error {
	kafkaResult, err := s.kafka.CommitMessage(msg)

	if err != nil {
		return err
	}

	log.Println("Success commit topic message", "topic", *kafkaResult[0].Topic, "result", kafkaResult[0].String())

	return nil
}

func getReportWhenIDNotFound(event *types.TraceEvent) *_report.Report {
	report := &_report.Report{
		ID:        event.Request.ID,
		Timestamp: getTimestamp(event.StartTime),
		TimeUnit:  _HOUR,
	}

	report.CallCountTotal, report.CallCountSuccess, report.CallCountFailed, report.CallCountOther, report.CallCountBlocked, report.ApiTimeTotal = calcCallCount(report, event)

	return report
}

func getReportWhenIDFounded(v *_report.Report, event *types.TraceEvent) *_report.Report {
	v.Timestamp = getTimestamp(event.StartTime)
	v.CallCountTotal, v.CallCountSuccess, v.CallCountFailed, v.CallCountOther, v.CallCountBlocked, v.ApiTimeTotal = calcCallCount(v, event)
	return v
}

// getCallCount
// Calculation of callCount through this document.
// https://learn.microsoft.com/en-us/rest/api/apimanagement/reports/list-by-api?view=rest-apimanagement-2024-05-01&tabs=HTTP#reportrecordcontract
func calcCallCount(v *_report.Report, event *types.TraceEvent) (callTotal, success, failed, other, blocked, apiTimeTotal int64) {
	// handleReport 해당 함수에서 Lock을 제어하고 있기 떄문에, Atomic은 굳이 필요가 없다.
	callTotal, success, failed, other, blocked, apiTimeTotal = v.CallCountTotal, v.CallCountSuccess, v.CallCountFailed, v.CallCountOther, v.CallCountBlocked, v.ApiTimeTotal

	callTotal++
	apiTimeTotal += int64(event.Latency)

	switch {
	case event.Response.IsCallFailed():
		failed++
	case event.Response.IsCallSuccessful():
		success++
	case event.Response.IsCallBlocked():
		blocked++
	default:
		other++
	}

	return callTotal, success, failed, other, blocked, apiTimeTotal
}

// getTimestamp
// remove minutes just get hour time
func getTimestamp(eventTime int64) int64 {
	t := time.Unix(0, eventTime*int64(time.Millisecond))
	truncatedTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	return truncatedTime.Unix()
}
