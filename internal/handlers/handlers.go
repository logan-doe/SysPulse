package handlers

import (
	"encoding/json"
	"net/http"
	"syspulse/internal/models"
	"syspulse/internal/services"
	"time"
)

var metricsService *services.MetricsService
var alertsService *services.AlertService

func SetMetricService(service *services.MetricsService) {
	metricsService = service
}

func SetAlertService(service *services.AlertService) {
	alertsService = service
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if alertsService == nil {
		http.Error(w, `{"error":"Alert service is not initialized"}`, http.StatusServiceUnavailable)
		return
	}

	history := alertsService.GetAlertHistory()
	json.NewEncoder(w).Encode(history)
}

func AlertConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if alertsService == nil {
		http.Error(w, `{"error": "Alert service not initialized"}`, http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case "GET":
		config := alertsService.GetConfig()
		json.NewEncoder(w).Encode(config)
	case "POST":
		var config models.AlertConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		alertsService.UpdateConfig(config)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func ClearAlertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if alertsService == nil {
		http.Error(w, `{"error":"alert service not initialised"}`, http.StatusServiceUnavailable)
		return
	}

	alertsService.ClearHistory()
	json.NewEncoder(w).Encode(map[string]string{"status": "cleared"})

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
