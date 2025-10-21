package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"syspulse/internal/config"
	"syspulse/internal/handlers"
	"syspulse/internal/services"
)

var (
	version        = "dev" // App version
	metricsService *services.MetricsService
	wsService      *services.WebSocketService
)

func main() {
	log.Printf("⚡️ SysPulse is started. Version (%s)\n", version)
	log.Printf("📊 Real Time System Monitor")
	log.Printf("🔌 Web Socket support enabled")

	// initialize services
	metricsService = services.NewMetricsService()
	wsService = services.NewWebSocketService()

	handlers.SetMetricService(metricsService)

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

	http.HandleFunc("/ws", wsService.HandleConnection)

	// ─── Main Page ───────────────────────────────────────────────────────
	http.HandleFunc("/", handlers.IndexHandler)

}

func startMetricBroadcast() { // starting periodically sending metrics
	ticker := time.NewTicker(2 * time.Second) // 2 seconds update rate
	defer ticker.Stop()

	for range ticker.C {
		metrics := metricsService.GetSystemMetrics()
		wsService.BroadcastMessage(metrics)

		clientsCount := wsService.GetConnectedClientCount()
		if clientsCount > 0 {
			log.Printf(" sending metrics to %d clients", clientsCount)
		}
	}

}
