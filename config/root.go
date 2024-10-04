package config

import (
	"github.com/naoina/toml"
	"os"
)

func NewConfig(tracePath, gatewayPath string) (*TraceCfg, *GatewayCfg) {
	t := new(TraceCfg)

	if f, err := os.Open(tracePath); err != nil {
		panic(err)
	} else {
		if err = toml.NewDecoder(f).Decode(t); err != nil {
			panic(err)
		}
	}

	gw := new(GatewayCfg)

	if f, err := os.Open(gatewayPath); err != nil {
		panic(err)
	} else {
		if err = toml.NewDecoder(f).Decode(gw); err != nil {
			panic(err)
		}
	}

	return t, gw
}
