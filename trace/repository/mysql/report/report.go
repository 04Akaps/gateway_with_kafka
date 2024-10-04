package report

import (
	"context"
	"database/sql"

	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql/utils"
	"github.com/upper/db/v4"

	. "github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql/report/query"
	_report "github.com/04Akaps/gateway_with_kafka.git/trace/types/report"
)

type report struct {
	report db.Collection
	utils  utils.SqlUtils
}

type ReportImpl interface {
	UpdateReportUsingBulkWithTx(tx db.Session, eventMapper map[string]*_report.Report) (sql.Result, error)
	TxContext(fn func(sess db.Session) error) error
}

var _ ReportImpl = (*report)(nil)

func NewReport(table db.Collection, utils utils.SqlUtils) ReportImpl {
	return &report{
		report: table,
		utils:  utils,
	}
}

func (s *report) UpdateReportUsingBulkWithTx(sess db.Session, mapper map[string]*_report.Report) (sql.Result, error) {
	query, err := ReportBulkQuery(mapper)
	if err != nil {
		return nil, err
	}

	result, err := sess.SQL().Exec(query)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *report) TxContext(fn func(sess db.Session) error) error {
	return s.session().TxContext(context.Background(), fn, nil)
}

func (s *report) session() db.Session {
	return s.report.Session()
}
