package config

type KafkaConsumerCfg struct {
	URI             string
	GroupID         string
	AutoOffsetReset string
	Polling         int
}

type TraceCfg struct {
	Kafka map[string]KafkaConsumerCfg

	MySQL map[string]struct {
		Database string
		Host     string
		User     string
		Password string

		Collections map[string]string
	}

	ServiceInfo struct {
		WorkChanLength int64
	}
}
