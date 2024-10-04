package gateway

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/global"
)

func NewGatewayAPI(cfg config.GatewayCfg) {

	global.InitializeWork.Done()
}
