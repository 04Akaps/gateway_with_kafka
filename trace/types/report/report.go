package report

//unique (id, timestamp)

type Report struct {
	ID               string `db:"id"`
	Timestamp        int64  `db:"timestamp"`
	TimeUnit         string `db:"time_unit"`
	ApiTimeTotal     int64  `db:"api_time_total"`
	CallCountTotal   int64  `db:"call_count_total"`
	CallCountSuccess int64  `db:"call_count_success"`
	CallCountFailed  int64  `db:"call_count_failed"`
	CallCountOther   int64  `db:"call_count_other"`
	CallCountBlocked int64  `db:"call_count_blocked"`
}
