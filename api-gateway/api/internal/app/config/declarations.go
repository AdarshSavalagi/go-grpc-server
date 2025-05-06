package config

// Config holds all the application configurations, including server, database, security, and middleware settings.
type Config struct {
	Server     TServerConfig     // Server-related configurations.
	Redis      TRedisConfig      // Redis-related configurations.
	Database   TDatabaseConfig   // Database-related configurations.
	Security   TSecurityConfig   // Security-related configurations.
	Middleware TMiddlewareConfig // Middleware-related configurations.
}

// TServerConfig defines the configuration for the server.
type TServerConfig struct {
	Environment string // The current environment (e.g., development, production).
	Port        int    // The port on which the server runs.
}

// TRedisConfig holds the configuration for Redis.
type TRedisConfig struct {
	Host     string // Redis server host.
	Port     int    // Redis server port.
	Password string // Redis authentication password (if required).
	DB       int    // Redis database number.
}

// TDatabaseConfig defines the configuration for the database connection.
type TDatabaseConfig struct {
	Host          string // Database host address.
	Port          int    // Database port number.
	User          string // Database username.
	Password      string // Database password.
	DBName        string // Database name.
	AutoMigration bool   // Whether to run automatic migrations.
	SSLMode       string // SSL mode for database connection.
}

// TSecurityConfig defines the security settings for authentication and authorization.
type TSecurityConfig struct {
	AllowedAuthStrategies  []string // List of allowed authentication strategies.
	AuthEnabled            bool     // Whether authentication is enabled.
	AuthMethod             string   // The authentication method in use.
	AllowedAuthzStrategies []string // List of allowed authorization strategies.
	AuthzMethod            string   // The authorization method in use.
}

// TMiddlewareConfig holds configuration flags for various middleware components.
type TMiddlewareConfig struct {
	ClientID           bool // Whether Client-ID middleware is enabled.
	DeviceID           bool // Whether Device-ID middleware is enabled.
	RequestID          bool // Whether Request-ID middleware is enabled.
	Idempotency        bool // Whether Idempotency middleware is enabled.
	GzipCompression    bool // Whether Gzip compression middleware is enabled.
	UserAgent          bool // Whether User-Agent middleware is enabled.
	SecurityHeaders    bool // Whether security headers middleware is enabled.
	CORS               bool // Whether CORS middleware is enabled.
	Cache              bool // Whether default cache behavior is enabled.
	ContentNegotiation bool // Whether content negotiation middleware is enabled.
	Referer            bool // Whether Referer header middleware is enabled.
	Cookies            bool // Whether Cookies middleware is enabled.
	RateLimit          bool // Whether rate-limiting middleware is enabled.
	Logger             bool // Whether logging middleware is enabled.
	Metrics            bool // Whether metrics middleware is enabled.
	Auth               bool // Whether authentication middleware is enabled.
	Profiler           bool // Whether profiling middleware is enabled.
}

// TSMTPConfig holds the configuration for SMTP email services.
type TSMTPConfig struct {
	SMTPHost        string // SMTP server host.
	SMTPPort        int    // SMTP server port.
	Username        string // SMTP username.
	Password        string // SMTP password.
	FromAddress     string // Default sender email address.
	FromName        string // Default sender name.
	TemplatePath    string // Path to email templates.
	DefaultTemplate string // Default email template name.
}

// Global variables holding configuration instances.
var (
	ServerConfig     TServerConfig     // Stores server configuration.
	RedisConfig      TRedisConfig      // Stores Redis configuration.
	DBConfig         TDatabaseConfig   // Stores database configuration.
	SecurityConfig   TSecurityConfig   // Stores security configuration.
	MiddlewareConfig TMiddlewareConfig // Stores middleware configuration.
)
