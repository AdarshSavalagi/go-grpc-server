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

	// Initialize server configuration.
	ServerConfig := TServerConfig{
		Port:        viper.GetInt("server.port"),
		Environment: *env,
		MaxRecSize: viper.GetInt("server.max-rec-size"),
		MaxSendSize: viper.GetInt("server.max-send-size"),
	}

	if ServerConfig.Port == 0 {
		return nil, errors.New("server port cannot be 0")
	}

	// Initialize Kafka config
	KafkaConfig := TKafkaConfig{
		Brokers:           viper.GetStringSlice("kafka.brokers"),
		RetryAttempts:     viper.GetInt("kafka.retry_attempts"),
		Topics:            viper.GetStringMapString("kafka.topics"),
		BufferChannelSize: viper.GetInt("kafka.buffer_channel_size"),
	}

	return &Config{
		Server: ServerConfig,
		Kafka:  KafkaConfig,
	}, nil
}
