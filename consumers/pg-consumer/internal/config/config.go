package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() (*Config, error) {
	env := flag.String("env", "development", "set the environment (e.g. development, production)")
	flag.Parse()

	// Default to production if the environment is not set.
	if *env == "" {
		*env = "production"
	}

	// Load the environment-specific .env file.
	err := godotenv.Load(".env." + *env)
	if err != nil {
		log.Fatalf("Error loading .env file for %s environment", *env)
	}

	viper.SetConfigName("config." + *env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".") // Also look in the root directory.

	// Read the configuration file.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	// Enable automatic environment variable binding.
	viper.AutomaticEnv()

	// Bind specific environment variables to configuration values.
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "DATABASE_NAME")
	viper.BindEnv("database.ssl-mode", "DATABASE_SSL_MODE")

	DatabaseConfig = TDatabaseConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.ssl-mode"),
		AutoMigration: viper.GetBool("database.auto-migration"),
	}

	KafkaConfig = TKafkaConfig{
		Brokers: viper.GetStringSlice("kafka.brokers"),
		RetryAttempts:    viper.GetInt("kafka.retry-attempts"),
		Topics:          viper.GetStringMapString("kafka.topics"),
		BufferChannelSize: viper.GetInt("kafka.buffer-channel-size"),
	}

	DatabaseConfig = TDatabaseConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.ssl-mode"),
	}

	return &Config{
		Kafka: KafkaConfig,
		Database: DatabaseConfig,
	}, nil
}
