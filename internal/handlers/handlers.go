package handlers

import (
	"encoding/json"
	"net/http"
	"syspulse/internal/services"
	"time"
)

var metricsService = services.NewMetricsService()

func SetMetricService(service *services.MetricsService) {
	metricsService = service
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service":   "SysPulse Monitor",
		"status":    "healthy",
		"timastamp": time.Now().UTC(),
		"features":  []string{"WebSocket", "Real-Time", "Metrics"},
	})
}

func VersionHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"service": "SysPulse Monitor",
			"version": version,
		})
	}
}

func ClientsHandler(wsService *services.WebSocketService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connected_clients": wsService.GetConnectedClientCount(),
			"timestamp":         time.Now().UTC(),
		})
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if metricsService == nil {
		http.Error(w, `{"error": "metrics service not initialized"}`, http.StatusServiceUnavailable)
		return
	}

	metrics := metricsService.GetSystemMetrics()
	json.NewEncoder(w).Encode(metrics)
}

/*
	 func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"service": "SysPulse Monitor",
			"message": "message from web socket",
			"status":  "🚧 under maintenance 🚧",
		})
	}
*/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}
