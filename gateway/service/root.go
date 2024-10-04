package service

import "github.com/04Akaps/gateway_with_kafka.git/config"

type service struct {
	cfg config.GatewayCfg
}

type ServiceImpl interface {
}

func NewService(cfg config.GatewayCfg) ServiceImpl {
	s := &service{cfg: cfg}

	return s
}
