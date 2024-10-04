package mysql

import (
	"database/sql"
	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql"
	_report "github.com/04Akaps/gateway_with_kafka.git/trace/types/report"
	"github.com/upper/db/v4"
)

type mysqlService struct {
	SQL *mysql.MySQL
}

type MySQLImpl interface {
	UpdateReportUsingBulkWithTx(tx db.Session, eventMapper map[string]*_report.Report) (sql.Result, error)
	TxContext(fn func(sess db.Session) error) error
}

var _ MySQLImpl = (*mysqlService)(nil)

func NewMySQLService(sql *mysql.MySQL) MySQLImpl {
	s := &mysqlService{
		SQL: sql,
	}

	return s
}

func (m *mysqlService) UpdateReportUsingBulkWithTx(tx db.Session, eventMapper map[string]*_report.Report) (sql.Result, error) {
	impl := m.SQL.GetSubscriptionReport()

	if res, err := impl.UpdateReportUsingBulkWithTx(tx, eventMapper); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// TxContext(fn func(sess db.Session) error) error
func (m *mysqlService) TxContext(fn func(sess db.Session) error) error {
	impl := m.SQL.GetSubscriptionReport()
	return impl.TxContext(fn)
}
