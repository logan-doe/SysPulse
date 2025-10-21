package services

import (
	"os"
	"runtime"
	"syspulse/internal/models"
	"time"
)

type MetricsService struct {
	lastCPUStats *CPUStats
}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

func (ms *MetricsService) GetSystemMetrics() models.SystemMetrics {
	return models.SystemMetrics{
		TimeStamp: time.Now(),
		CPU:       ms.getCPUInfo(),
		Memory:    ms.getMemoryInfo(),
		Disk:      ms.getDiskInfo(),
		System:    ms.getSystemInfo(),
	}
}

func (ms *MetricsService) getCPUInfo() models.CPUInfo {
	return models.CPUInfo{
		Usage:  ms.calculateCPUUsage(), // default random value
		Cores:  runtime.NumCPU(),
		Load1:  1.2,
		Load5:  1.5,
		Load15: 1.3,
	}
}

func (ms *MetricsService) getMemoryInfo() models.MemInfo {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	total := memStats.Sys
	used := memStats.Alloc
	available := total - used
	usage := float64(used) / float64(total) * 100

	return models.MemInfo{
		Total:     total,
		Used:      used,
		Available: available,
		Usage:     usage,
	}
}

func (ms *MetricsService) getDiskInfo() models.DiskInfo {
	return models.DiskInfo{
		Total: 2000 * 1024 * 1024 * 1024,
		Used:  460 * 1024 * 1024 * 1024,
		Free:  1540 * 1024 * 1024 * 1024,
		Usage: 35.0,
	}
}

func (ms *MetricsService) getSystemInfo() models.SystemInfo {

	hostname, _ := os.Hostname()
	return models.SystemInfo{
		Hostname: hostname,
		OS:       runtime.GOOS,
		Platform: runtime.GOARCH,
		Uptime:   ms.getSystemUpTime(),
	}
}

type CPUStats struct {
	Time  time.Time
	Usage float64
}

func (ms *MetricsService) calculateCPUUsage() float64 {
	// Более реалистичная симуляция загрузки CPU
	// Основана на времени и случайных факторах
	now := time.Now()
	second := now.Second()

	// Симулируем разную загрузку в зависимости от времени
	var usage float64

	switch {
	case second < 15:
		usage = 10.0 + float64(second) // Плавный рост 10-25%
	case second < 30:
		usage = 40.0 - float64(second-15) // Плавное снижение 40-25%
	case second < 45:
		usage = 15.0 + float64(second-30)*0.5 // Медленный рост 15-22.5%
	default:
		usage = 60.0 - float64(second-45) // Снижение 60-45%
	}

	// Добавляем небольшую случайную вариацию
	if second%2 == 0 {
		usage += 2.0
	} else {
		usage -= 1.5
	}

	// Ограничиваем диапазон
	if usage < 5.0 {
		usage = 5.0
	}
	if usage > 95.0 {
		usage = 95.0
	}

	return usage
}

func (ms *MetricsService) getSystemUpTime() uint64 {
	baseUpTime := uint64(24 * 60 * 60)
	additional := uint64(time.Now().Minute() * 60)
	return baseUpTime + additional
}
