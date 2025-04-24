package config

import (
	"errors"
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitConfig() (*Config, error) {
	env := flag.String("env", "development", "set the environment (e.g. development, production)")
	flag.Parse()

	log.Printf("Environment: %s", *env)

	// Default to production if the environment is not set.
	if *env == "" {
		*env = "production"
	}

	// Load environment variables from the .env file.
	err := godotenv.Load(".env." + *env)
	if err != nil {
		log.Fatalf("Error loading .env file for %s environment", *env)
	}

	// Set up viper for reading YAML config.
	viper.SetConfigName("config." + *env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Automatically map environment variables to config
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.refresh_token_secret", "JWT_REFRESH_TOKEN_SECRET")
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("kafka.sasl_user", "KAFKA_SASL_USER")
	viper.BindEnv("kafka.sasl_password", "KAFKA_SASL_PASSWORD")
	viper.BindEnv("kafka.sasl_mechanism", "KAFKA_SASL_MECHANISM")

	// Initialize server configuration.
	ServerConfig := TServerConfig{
		Port:        viper.GetInt("server.port"),
		Environment: *env,
	}

	if ServerConfig.Port == 0 {
		return nil, errors.New("server port cannot be 0")
	}

	// Initialize Kafka config
	KafkaConfig := TKafkaConfig{
		Brokers:       viper.GetStringSlice("kafka.brokers"),
		ClientID:      viper.GetString("kafka.client_id"),
		Acks:          viper.GetString("kafka.acks"),
		Async:         viper.GetBool("kafka.async"),
		RetryAttempts: viper.GetInt("kafka.retry_attempts"),
		WriteTimeout:  viper.GetInt("kafka.write_timeout"),
		ReadTimeout:   viper.GetInt("kafka.read_timeout"),
		Topics: map[string]string{
			"logs":    viper.GetString("kafka.topics.logs"),
			"events":  viper.GetString("kafka.topics.events"),
			"context": viper.GetString("kafka.topics.context"),
		},
		EnableTLS:     viper.GetBool("kafka.enable_tls"),
		EnableSASL:    viper.GetBool("kafka.enable_sasl"),
		SASLUser:      viper.GetString("kafka.sasl_user"),
		SASLPassword:  viper.GetString("kafka.sasl_password"),
		SASLMechanism: viper.GetString("kafka.sasl_mechanism"),
	}

	return &Config{
		Server: ServerConfig,
		Kafka:  KafkaConfig,
	}, nil
}
