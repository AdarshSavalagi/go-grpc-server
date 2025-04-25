package config

type Config struct {
	Server TServerConfig
	Kafka TKafkaConfig
}

type TServerConfig struct {
	Environment string
	Port 	 int
}



type TKafkaConfig struct {
	Brokers       []string          `mapstructure:"brokers"`
	RetryAttempts int               `mapstructure:"retry_attempts"`

	Topics map[string]string `mapstructure:"topics"` // "logs", "events", "context"
	BufferChannelSize int `mapstructure:"buffer_channel_size"`
}

// Global variables holding configuration instances.
var (
	ServerConfig     TServerConfig     // Stores server configuration.
	KafkaConfig      TKafkaConfig      // Stores Kafka configuration.
)