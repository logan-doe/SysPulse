package services

import (
	"fmt"
	"log"
	"sync"
	"syspulse/internal/models"
	"time"
)

type AlertService struct {
	mu           sync.Mutex
	alerts       []models.Alert
	maxAlerts    int
	config       models.AlertConfig
	activeAlerts map[string]bool
}

func NewAlertService() *AlertService {
	return &AlertService{
		alerts:       make([]models.Alert, 0),
		maxAlerts:    50,
		activeAlerts: make(map[string]bool),
		config: models.AlertConfig{
			CPUTreshold:  75.0,
			RAMTreshold:  75.0,
			DiskTreshold: 85.0,
			Enabled:      true,
		},
	}
}

func (as *AlertService) createAlert(alertType string, value, threshold float64, timestamp time.Time) models.Alert {
	level := "warning"
	if value >= threshold+10 {
		level = "critical"
	}

	alert := models.Alert{
		ID:        generateID(alertType, timestamp),
		Type:      alertType,
		Message:   as.generateAlertMessage(alertType, value, threshold, level),
		Level:     level,
		Threshold: threshold,
		Timestamp: timestamp,
		Active:    true,
	}

	as.activeAlerts[alertType] = true
	return alert
}

func generateID(alertType string, timestamp time.Time) string {
	return fmt.Sprintf("%s-%d", alertType, timestamp.Unix())
}

func (as *AlertService) CheckMetrics(metrics models.SystemMetrics) []models.Alert {
	if !as.config.Enabled {
		return []models.Alert{}
	}

	var newAlerts []models.Alert
	now := time.Now()

	if metrics.CPU.Usage >= as.config.CPUTreshold {
		alert := as.createAlert("CPU", metrics.CPU.Usage, as.config.CPUTreshold, now)
		newAlerts = append(newAlerts, alert)
	} else {
		as.resolveAlert("CPU")
	}

	if metrics.Memory.Usage >= as.config.RAMTreshold {
		alert := as.createAlert("RAM", metrics.Memory.Usage, as.config.RAMTreshold, now)
		newAlerts = append(newAlerts, alert)
	} else {
		as.resolveAlert("RAM")
	}

	if metrics.Disk.Usage >= as.config.DiskTreshold {
		alert := as.createAlert("DISK", metrics.Disk.Usage, as.config.DiskTreshold, now)
		newAlerts = append(newAlerts, alert)
	} else {
		as.resolveAlert("DISK")
	}

	if len(newAlerts) > 0 {
		as.addAlert(newAlerts)
	}

	return newAlerts
}

func (as *AlertService) generateAlertMessage(alertType string, value, threshold float64, level string) string {
	switch alertType {
	case "CPU":
		return formatAlertMessage("CPU", value, threshold, level, "Loading")
	case "RAM":
		return formatAlertMessage("RAM", value, threshold, level, "Usage")
	case "DISK":
		return formatAlertMessage("DISK", value, threshold, level, "Usage")
	default:
		return formatAlertMessage("System", value, threshold, level, "Loading")

	}
}

func formatAlertMessage(resource string, value, threshold float64, level, metric string) string {
	levelTxt := "warning"
	if level == "critical" {
		levelTxt = level
	}

	return fmt.Sprintf("[%s] %s: %s %.1f%% (Threshold: %.1f%%)\n", resource, levelTxt, metric, value, threshold)
}

func (as *AlertService) resolveAlert(alertType string) { // removing alert when metric is getting under threshold
	if as.activeAlerts[alertType] {
		log.Printf("✅ Alert is removed")
		delete(as.activeAlerts, alertType)
	}
}

func (as *AlertService) addAlert(alerts []models.Alert) { // add alerts to history block
	as.mu.Lock()
	defer as.mu.Unlock()

	for _, alert := range alerts {
		as.alerts = append(as.alerts, alert)
		log.Printf("⚠️ New alert: %s\n", alert.Message)
	}

	if len(as.alerts) > as.maxAlerts { // if length of alerts more than max lenth of history - using FIFO
		as.alerts = as.alerts[len(as.alerts)-as.maxAlerts:]
	}
}

func (as *AlertService) GetAlertHistory() models.AlertHistory {
	as.mu.Lock()
	defer as.mu.Unlock()

	var history models.AlertHistory
	history.Alerts = make([]models.Alert, len(as.alerts))
	copy(history.Alerts, as.alerts)

	today := time.Now().Truncate(24 * time.Hour)

	for _, alert := range as.alerts {
		history.Stats.TotalAlerts++
		if alert.Active {
			history.Stats.ActiveAlerts++
		}
		if alert.Timestamp.After(today) {
			history.Stats.TodayAlerts++
		}
	}
	return history
}

func (as *AlertService) UpdateConfig(config models.AlertConfig) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.config = config
	log.Printf("🪪 Alert config updated: CPU = %.1f%%, RAM = %.1f%%, Disk = %.1f%%, Enabled = %v\n", config.CPUTreshold, config.RAMTreshold, config.DiskTreshold, config.Enabled)
}

func (as *AlertService) GetConfig() models.AlertConfig {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.config
}

func (as *AlertService) ClearHistory() {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.alerts = make([]models.Alert, 0)
	as.activeAlerts = make(map[string]bool)

	log.Printf("✅ Alert history is clear")
}
