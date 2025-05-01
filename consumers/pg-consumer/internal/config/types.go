package config


type Config struct {
	Kafka    TKafkaConfig
	Database TDatabaseConfig
}

type TKafkaConfig struct {
	Brokers           []string
	RetryAttempts     int
	Topics            map[string]string
	BufferChannelSize int
}

type TDatabaseConfig struct {
	Host          string
	Port          string
	User          string
	Password      string
	DBName        string
	AutoMigration bool
	SSLMode       string
}

func (d *TDatabaseConfig) DSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.DBName +
		" port=" + d.Port +
		" sslmode=" + d.SSLMode
}

var (
	KafkaConfig    TKafkaConfig
	DatabaseConfig TDatabaseConfig
)
