package trace

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/trace/app"
)

func NewTraceModule(cfg config.TraceCfg) {
	app := app.NewApp(cfg)
	app.Run()
}
