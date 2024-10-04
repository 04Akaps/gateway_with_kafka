package utils

import (
	"github.com/upper/db/v4"
	"log"
)

type SqlUtils struct{}

func NewSqlUtils() SqlUtils {
	return SqlUtils{}
}

func (s *SqlUtils) IterCloseHandler(iter db.Iterator) {
	if err := iter.Close(); err != nil {
		log.Println("Failed to close iterator handler", "err", err)
	}
}
