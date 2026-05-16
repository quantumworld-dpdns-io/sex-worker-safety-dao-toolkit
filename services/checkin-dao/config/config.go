package config

import "os"

type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	QdrantHost        string
	QdrantPort        string
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioFromNumber  string
	EmergencyWebhook  string
}

func Load() *Config {
	cfg := &Config{
		Port:             getEnv("PORT", "3002"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		QdrantHost:       getEnv("QDRANT_HOST", "localhost"),
		QdrantPort:       getEnv("QDRANT_PORT", "6333"),
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
		EmergencyWebhook: getEnv("EMERGENCY_WEBHOOK_URL", ""),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
