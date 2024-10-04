package main

import (
	"flag"
	"fmt"
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/gateway"
	"github.com/04Akaps/gateway_with_kafka.git/global"
	"github.com/04Akaps/gateway_with_kafka.git/trace"
)

var (
	traceCfg   = flag.String("trace-cfg", "./trace/config.toml", "trace cfg toml file path")
	gatewayCfg = flag.String("gateway-cfg", "./gateway/config.toml", "gateway cfg toml file path")
)

func init() {
	flag.Parse()
	global.InitializeWork.Add(2)
}

func main() {
	tcfg, gwcfg := config.NewConfig(*traceCfg, *gatewayCfg)

	go trace.NewTraceModule(*tcfg)
	go gateway.NewGatewayAPI(*gwcfg)

	global.InitializeWork.Wait()

	fmt.Println("여기 온다.")
	for {
	}
}
