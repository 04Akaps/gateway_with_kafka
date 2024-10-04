package mysql

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql/report"
	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql/utils"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

// For Remove Full Group By
const sqlMode = "'STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ALLOW_INVALID_DATES,ERROR_FOR_DIVISION_BY_ZERO'"

type MySQL struct {
	sess db.Session

	cfg config.TraceCfg

	subscriptionReport report.ReportImpl
	utils              utils.SqlUtils
}

func NewMySQL(config config.TraceCfg) (*MySQL, error) {
	sql := &MySQL{
		cfg: config,
	}

	var err error
	cfg := config.MySQL["link"]

	if sql.sess, err = mysql.Open(mysql.ConnectionURL{
		Database: cfg.Database,
		Host:     cfg.Host,
		User:     cfg.User,
		Password: cfg.Password,
		// -> if need
		//Options: map[string]string{
		//	"sql_mode": sqlMode,
		//},
	}); err != nil {
		return nil, err
	} else if err = sql.sess.Ping(); err != nil {
		return nil, err
	} else {
		sql.utils = utils.NewSqlUtils()

		sql.subscriptionReport = report.NewReport(sql.sess.Collection(cfg.Collections["Report"]), sql.utils)
	}

	return sql, nil
}

func (m *MySQL) GetSubscriptionReport() report.ReportImpl {
	return m.subscriptionReport
}
