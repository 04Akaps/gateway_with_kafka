package query

import (
	"fmt"
	"github.com/04Akaps/gateway_with_kafka.git/trace/types/report"
)

const (
	seReportBulk = "" +
		"INSERT INTO report (%s) VALUES%s " +
		"ON DUPLICATE KEY UPDATE " +
		"api_time_total = api_time_total + VALUES(api_time_total)," +
		"call_count_total = call_count_total + VALUES(call_count_total)," +
		"call_count_success = call_count_success + VALUES(call_count_success)," +
		"call_count_failed = call_count_failed + VALUES(call_count_failed)," +
		"call_count_other = call_count_other + VALUES(call_count_other)," +
		"call_count_blocked = call_count_blocked + VALUES(call_count_blocked);"

	setReportBulkField = "" +
		"id, timestamp, time_unit, api_time_total, call_count_total, call_count_success, call_count_failed, call_count_other, call_count_blocked"
)

func ReportBulkQuery(eventMapper map[string]*report.Report) (string, error) {

	var values string

	for _, rp := range eventMapper {
		values += fmt.Sprintf(
			" (%d, '%s', %d, %d, %d, %d, %d, %d),",
			rp.Timestamp, rp.TimeUnit, rp.ApiTimeTotal,
			rp.CallCountTotal, rp.CallCountSuccess, rp.CallCountFailed, rp.CallCountOther, rp.CallCountBlocked,
		)
	}

	// remove last comma
	if len(values) > 0 {
		values = values[:len(values)-1]
	}

	return fmt.Sprintf(seReportBulk, setReportBulkField, values), nil
}
