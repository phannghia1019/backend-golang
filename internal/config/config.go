package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	NodeEnv              string
	Port                 string
	JWTSecretKey         string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURI    string
	GoogleAESKey         string
	GoogleErrorPermURL   string
	MongoUser            string
	MongoPass            string
	MongoHost            string
	MongoPort            string
	MongoDB              string
	SAPEndpoint          string
	SAPDBHost            string
	SAPDBUser            string
	SAPDBPass            string
	SAPDBPort            string
	SAPDBList            string
	HostURL              string
	AppManagerAPIHost    string
	LogLevel             string
}

var cfg *Config

// Load loads configuration from environment
func Load(envPath string) (*Config, error) {
	// Load .env file if exists
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			// Not an error if file doesn't exist
		}
	}

	cfg = &Config{
		NodeEnv:              getEnv("NODE_ENV", "development"),
		Port:                 getEnv("PORT", "3000"),
		JWTSecretKey:         getEnv("JWT_SECRET_KEY", ""),
		GoogleClientID:       getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:   getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:    getEnv("GOOGLE_REDIRECT_URI", ""),
		GoogleAESKey:         getEnv("GOOGLE_AES_KEY", ""),
		GoogleErrorPermURL:   getEnv("GOOGLE_ERROR_PERMISSION_URL", ""),
		MongoUser:            getEnv("MONGO_USER", ""),
		MongoPass:            getEnv("MONGO_PASS", ""),
		MongoHost:            getEnv("MONGO_HOST", "localhost"),
		MongoPort:            getEnv("MONGO_PORT", "27017"),
		MongoDB:              getEnv("MONGO_DB", ""),
		SAPEndpoint:          getEnv("SAP_ENDPOINT", ""),
		SAPDBHost:            getEnv("SAP_DB_HOST", ""),
		SAPDBUser:            getEnv("SAP_DB_USER", ""),
		SAPDBPass:            getEnv("SAP_DB_PASS", ""),
		SAPDBPort:            getEnv("SAP_DB_PORT", "1433"),
		SAPDBList:            getEnv("SAP_DB_LIST", "HBCCORP,KBI"),
		HostURL:              getEnv("HOST_URL", ""),
		AppManagerAPIHost:    getEnv("APP_MANAGER_API_HOST", ""),
		LogLevel:             getEnv("LOG_LEVEL", "INFO"),
	}

	return cfg, nil
}

// Get returns the current configuration
func Get() *Config {
	return cfg
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
