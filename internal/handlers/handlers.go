package handlers

import (
	"encoding/json"
	"net/http"
	"syspulse/internal/models"
	"time"
)

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

	// ─── Main Metrics Handler ────────────────────────────────────────────
	metrics := models.SystemMetrics{
		TimeStamp: time.Now().UTC(),
		CPU: models.CPUInfo{
			Usage:  10.0,
			Cores:  12,
			Load1:  1.5,
			Load5:  2.0,
			Load15: 1.8,
		},
		Memory: models.MemInfo{
			Total:     17179869184, // 16GB
			Used:      4294967296,  // 4GB
			Available: 12884901888, // 12GB
			Usage:     25.0,
		},
		Disk: models.DiskInfo{
			Total: 536870912000, // 500GB
			Used:  161061273600, // 150GB
			Free:  375809638400, // 350GB
			Usage: 30.0,
		},
		System: models.SystemInfo{
			Hostname: "default Name",
			OS:       "test OS",
			Platform: "x64",
			Uptime:   86400,
		},
	}

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
