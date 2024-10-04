package error

import (
	"errors"
	"strings"
)

type ErrString string

const (
	ErrConsumerNotFound  = ErrString("consumer can't find")
	ErrTopicNotSupported = ErrString("topic not supported")
)

func (e ErrString) toString() string {
	return string(e)
}

func New(code ErrString, msg ...string) error {

	text := code.toString()

	if len(msg) > 0 {
		text += ": " + strings.Join(msg, " ") // 메시지를 공백으로 결합
	}

	return errors.New(text)
}
