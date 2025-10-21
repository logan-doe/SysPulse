package main

import (
	"fmt"
	"log"
	"net/http"

	"syspulse/internal/config"
	"syspulse/internal/handlers"
)

var version = "dev" // App version

func main() {
	log.Printf("⚡️ SysPulse is started. Version (%s)\n", version)
	log.Printf("📊 Real Time System Monitor")

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
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	// ─── Main Page ───────────────────────────────────────────────────────
	http.HandleFunc("/", handlers.IndexHandler)

}
