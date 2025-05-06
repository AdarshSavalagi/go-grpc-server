// Contents: Configuration settings for the application.
package config

import (
	"api/internal/utils/auth"
	"errors"
	"flag"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// ConfigInit initializes the application configuration by loading environment variables
// and configuration files. It sets up JWT authentication, database, security, middleware,
// and Redis configurations.
//
// Returns:
//   - *Config: A pointer to the initialized Config structure.
//   - error: An error if the configuration fails to load.
func ConfigInit() (*Config, error) {

	// Parse the environment argument from command-line flags.
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

	// Configure Viper to read YAML configuration files.
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

	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.refresh-secret", "JWT_REFRESH_SECRET")

	// Load JWT secret key from config or fail if not set.
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		log.Fatalf("No JWT secret found in config.")
	}

	// Set JWT secrets in the authentication module.
	viper.Set("jwt.secret", secret)

	refreshSecret := viper.GetString("jwt.refresh-secret")
	if refreshSecret == "" {
		log.Fatalf("No JWT Refresh secret found in config.")
	}

	viper.Set("jwt.refresh-secret", refreshSecret)

	// Configure JWT authentication settings.
	auth.SetJwtKey([]byte(secret))
	auth.SetJwtRefreshKey([]byte(refreshSecret))
	auth.SetAccessTokenExpiry(viper.GetInt("jwt.access-token-expiry"))
	auth.SetRefreshTokenExpiry(viper.GetInt("jwt.refresh-token-expiry"))

	// Initialize server configuration.
	ServerConfig = TServerConfig{
		Port:        viper.GetInt("server.port"),
		Environment: *env,
	}

	if ServerConfig.Port == 0 {
		return nil, errors.New("server port cannot be 0")
	}

	// Initialize database configuration.
	DBConfig = TDatabaseConfig{
		Host:          viper.GetString("database.host"),
		Port:          viper.GetInt("database.port"),
		User:          viper.GetString("database.user"),
		Password:      viper.GetString("database.password"),
		DBName:        viper.GetString("database.dbname"),
		SSLMode:       viper.GetString("database.ssl-mode"),
		AutoMigration: viper.GetBool("db-migration"),
	}

	// Initialize security configuration.
	SecurityConfig = TSecurityConfig{
		AllowedAuthStrategies:  viper.GetStringSlice("security.allowed_authentication_strategies"),
		AuthEnabled:            viper.GetBool("security.authentication_enabled"),
		AuthMethod:             viper.GetString("security.authentication_method"),
		AllowedAuthzStrategies: viper.GetStringSlice("security.allowed_authorization_strategies"),
		AuthzMethod:            viper.GetString("security.authorization_method"),
	}

	// Initialize middleware configuration.
	MiddlewareConfig = TMiddlewareConfig{
		ClientID:           viper.GetBool("middlewares.client-id"),
		DeviceID:           viper.GetBool("middlewares.device-id"),
		RequestID:          viper.GetBool("middlewares.request-id"),
		Idempotency:        viper.GetBool("middlewares.idempotency"),
		GzipCompression:    viper.GetBool("middlewares.gzip-compression"),
		UserAgent:          viper.GetBool("middlewares.user-agent"),
		SecurityHeaders:    viper.GetBool("middlewares.security-headers"),
		CORS:               viper.GetBool("middlewares.cors"),
		Cache:              viper.GetBool("middlewares.default-cache-behavior"),
		ContentNegotiation: viper.GetBool("middlewares.content-negotiation"),
		Referer:            viper.GetBool("middlewares.referer"),
		Cookies:            viper.GetBool("middlewares.cookies"),
		RateLimit:          viper.GetBool("middlewares.rate-limit"),
		Logger:             viper.GetBool("middlewares.logger"),
		Metrics:            viper.GetBool("middlewares.metrics"),
		Profiler:           viper.GetBool("middlewares.profiler"),
		Auth:               viper.GetBool("security.authentication_enabled"),
	}

	// Bind SMTP environment variables.
	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.username", "SMTP_USERNAME")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")

	// Initialize Redis configuration.
	RedisConfig = TRedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	// Return the compiled configuration.
	return &Config{
		ServerConfig,
		RedisConfig,
		DBConfig,
		SecurityConfig,
		MiddlewareConfig,
	}, nil
}

// DSN constructs the Data Source Name for Redis connection.
//
// Returns:
//   - string: A properly formatted Redis DSN.
func (d *TRedisConfig) DSN() string {
	return d.Host + ":" + strconv.Itoa(d.Port)
}

// DSN constructs the Data Source Name for PostgreSQL connection.
//
// Returns:
//   - string: A properly formatted PostgreSQL DSN.
func (d *TDatabaseConfig) DSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.DBName +
		" port=" + strconv.Itoa(d.Port) +
		" sslmode=" + d.SSLMode
}
