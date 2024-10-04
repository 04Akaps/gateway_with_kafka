package config

type GatewayCfg struct {
	MySQL map[string]struct {
		Database string
		Host     string
		User     string
		Password string

		Collections map[string]string
	}

	Redis struct {
		DataSource string
		DB         int
		Password   string
		UserName   string
	}

	ServiceInfo struct {
		Port string
	}
}
