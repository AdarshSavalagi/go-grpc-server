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
	ClientID      string            `mapstructure:"client_id"`
	Acks          string            `mapstructure:"acks"`
	Async         bool              `mapstructure:"async"`
	RetryAttempts int               `mapstructure:"retry_attempts"`
	WriteTimeout  int               `mapstructure:"write_timeout"`
	ReadTimeout   int               `mapstructure:"read_timeout"`

	Topics map[string]string `mapstructure:"topics"` // "logs", "events", "context"

	// Auth (optional)
	EnableTLS     bool   `mapstructure:"enable_tls"`
	EnableSASL    bool   `mapstructure:"enable_sasl"`
	SASLUser      string `mapstructure:"sasl_user"`
	SASLPassword  string `mapstructure:"sasl_password"`
	SASLMechanism string `mapstructure:"sasl_mechanism"`
}

// Global variables holding configuration instances.
var (
	ServerConfig     TServerConfig     // Stores server configuration.
	KafkaConfig      TKafkaConfig      // Stores Kafka configuration.
)