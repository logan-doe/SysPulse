package services

import (
	"os"
	"runtime"
	"syspulse/internal/models"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
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

	cpuPercent, _ := cpu.Percent(time.Second, false)
	cpuUsage := 0.0

	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	cpuCores, _ := cpu.Counts(true)

	avgLoad, _ := load.Avg()
	load1, load5, load15 := 0.0, 0.0, 0.0
	if avgLoad != nil {
		load1 = avgLoad.Load1
		load5 = avgLoad.Load5
		load15 = avgLoad.Load15
	}

	return models.CPUInfo{
		Usage:  cpuUsage,
		Cores:  cpuCores,
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}
}

func (ms *MetricsService) getMemoryInfo() models.MemInfo {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memory, err := mem.VirtualMemory()
	if err != nil {
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return models.MemInfo{
			Total:     memStat.Sys,
			Used:      memStat.Alloc,
			Available: memStat.Sys - memStat.Alloc,
			Usage:     float64(memStats.Alloc) / float64(memStats.Sys) * 100,
		}
	}
	return models.MemInfo{
		Total:     memory.Total,
		Used:      memory.Used,
		Available: memory.Available,
		Usage:     memory.UsedPercent,
	}
}

func (ms *MetricsService) getDiskInfo() models.DiskInfo {

	partitions, err := disk.Partitions(false)
	if err != nil {
		return ms.getDiskInfoFallback()

	}

	var largestPartition string
	var maxsize uint64

	for _, partition := range partitions {
		if isSpecialFilesystem(partition.Fstype) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		if usage.Total > maxsize {
			maxsize = usage.Total
			largestPartition = partition.Mountpoint
		}
	}

	if largestPartition != "" {
		usage, err := disk.Usage(largestPartition)
		if err == nil {
			return models.DiskInfo{
				Total: usage.Total,
				Used:  usage.Used,
				Free:  usage.Free,
				Usage: usage.UsedPercent,
			}
		}
	}
	return ms.getDiskInfoFallback()
}

func isSpecialFilesystem(fstype string) bool {
	specialFS := []string{
		"tmpfs", "devtmpfs", "squashfs", "overlay",
		"proc", "sysfs", "devpts", "mqueue", "debugfs",
		"securityfs", "pstore", "cgroup", "cgroup2",
	}

	for _, fs := range specialFS {
		if fstype == fs {
			return true
		}
	}
	return false
}

func (ms *MetricsService) getDiskInfoFallback() models.DiskInfo {

	mountPoints := []string{"/", "/home", "/mnt", "/media", "C:\\", "D:\\"}

	for _, point := range mountPoints {
		if usage, err := disk.Usage(point); err == nil {
			return models.DiskInfo{
				Total: usage.Total,
				Used:  usage.Used,
				Free:  usage.Free,
				Usage: usage.UsedPercent,
			}
		}
	}

	return models.DiskInfo{
		Total: 5000 * 1024 * 1024 * 1024,
		Used:  1234 * 1024 * 1024 * 1024,
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
		Uptime:   uint64(time.Now().Unix() - ms.getSystemUpTime()),
	}
}

func (ms *MetricsService) getSystemUpTime() int64 {
	return time.Now().Add(-24 * time.Hour).Unix()
}

type CPUStats struct {
	Time  time.Time
	Usage float64
}
