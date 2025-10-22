package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"syspulse/internal/config"
	"syspulse/internal/handlers"
	"syspulse/internal/services"
)

var (
	version        = "1.0.0" // App version
	metricsService *services.MetricsService
	wsService      *services.WebSocketService
	alertService   *services.AlertService
)

func main() {
	log.Printf("⚡️ SysPulse is started. Version (%s)\n", version)
	log.Printf("📊 Real Time System Monitor")
	log.Printf("🔌 Web Socket support enabled")

	// initialize services
	metricsService = services.NewMetricsService()
	wsService = services.NewWebSocketService()
	alertService = services.NewAlertService()

	handlers.SetMetricService(metricsService)
	handlers.SetAlertService(alertService)

	// web socket service start
	wsService.Start()

	//sending out metrics
	go startMetricBroadcast()

	// loading config
	cfg := config.Load()

	setupRoutes()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🌐 server is running on http://localhost%s\n", addr)
	log.Printf("📱 mode: %s\n", cfg.Environment)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}

func setupRoutes() {

	// ─── Static Files ────────────────────────────────────────────────────
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ─── Handlers ────────────────────────────────────────────────────────
	http.HandleFunc("/api/version", handlers.VersionHandler(version))
	http.HandleFunc("/api/health", handlers.HealthHandler)
	http.HandleFunc("/api/metrics", handlers.MetricsHandler)
	http.HandleFunc("/api/clients", handlers.ClientsHandler(wsService))
	http.HandleFunc("/api/alerts/history", handlers.AlertHandler)
	http.HandleFunc("/api/alerts/config", handlers.AlertConfigHandler)
	http.HandleFunc("/api/alerts/clear", handlers.ClearAlertHandler)

	http.HandleFunc("/api/debug", func(w http.ResponseWriter, r *http.Request) {
		metrics := metricsService.GetSystemMetrics()
		json.NewEncoder(w).Encode(metrics)
	})

	http.HandleFunc("/ws", wsService.HandleConnection)

	// ─── Main Page ───────────────────────────────────────────────────────
	http.HandleFunc("/", handlers.IndexHandler)

}

func startMetricBroadcast() { // starting periodically sending metrics
	cfg := config.Load()

	ticker := time.NewTicker(time.Duration(cfg.UpdateInterval) * time.Millisecond) // update rate
	defer ticker.Stop()

	for range ticker.C {
		metrics := metricsService.GetSystemMetrics()

		alerts := alertService.CheckMetrics(metrics)
		metrics.Alerts = alerts

		wsService.BroadcastMessage(metrics)

		clientsCount := wsService.GetConnectedClientCount()
		if clientsCount > 0 && time.Now().Second()%10 == 0 {
			log.Printf("🪪 sending metrics to %d clients", clientsCount)
		}
	}

}
