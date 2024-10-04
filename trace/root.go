package trace

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/global"
)

func NewTraceModule(cfg config.TraceCfg) {

	global.InitializeWork.Done()
}
