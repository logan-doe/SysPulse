package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             string
	Environment      string
	UpdateInterval   int
	AlertThreshholds AlertConfig
}

type AlertConfig struct {
	CPU  float64
	RAM  float64
	Disk float64
}

func Load() *Config {
	cfg := &Config{
		Port:           getEnv("SYS_PULSE_PORT", "8080"),
		Environment:    getEnv("SYS_PULSE_ENVIRONMENT", "development"),
		UpdateInterval: getEnvInt("SYS_PULSE_UPDATE_INTERVAL", 500),
		AlertThreshholds: AlertConfig{
			CPU:  getEnvFloat("SYS_PULSE_ALERT_CPU", 80.0),
			RAM:  getEnvFloat("SYS_PULSE_ALERT_RAM", 85.0),
			Disk: getEnvFloat("SYS_PULSE_ALERT_DISK", 90.0),
		},
	}
	return cfg
}

func getEnv(key, defaultStr string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultStr
}

func getEnvInt(key string, defaultInt int) int {
	if value := os.Getenv(key); value != "" {
		valueInt, err := strconv.Atoi(value)
		if err == nil {
			return valueInt
		}
	}
	return defaultInt
}

func getEnvFloat(key string, defaultFloat float64) float64 {
	if value := os.Getenv(key); value != "" {
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return valueFloat
		}
	}
	return defaultFloat
}
