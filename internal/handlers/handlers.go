package handlers

import (
	"encoding/json"
	"net/http"
	"syspulse/internal/services"
	"time"
)

var metricsService = services.NewMetricsService()

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service":   "SysPulse Monitor",
		"status":    "healthy",
		"timastamp": time.Now().UTC(),
	})
}

func VersionHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service": "SysPulse Monitor",
			"version": version,
		})
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metrics := metricsService.GetSystemMetrics()
	json.NewEncoder(w).Encode(metrics)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"service": "SysPulse Monitor",
		"message": "message from web socket",
		"status":  "🚧 under maintenance 🚧",
	})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}
