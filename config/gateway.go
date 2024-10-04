package config

type GatewayCfg struct {
	MySQL map[string]struct {
		Database string
		Host     string
		User     string
		Password string

		Collections map[string]string
	}

	ServiceInfo struct {
		ServiceID      string
		WorkChanLength int64
		Port           string
	}
}
