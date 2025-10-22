package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"syspulse/internal/models"
	"time"

	"net"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	gnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type MetricsService struct {
	lastCPUStats       *CPUStats
	prevNetCounters    map[string]gnet.IOCountersStat // contains value of prev net counters
	prevNetTime        time.Time                      //time of prev time measure
	pingServers        []string                       // primary servers for measuring ping
	currentServerIndex int
	totalUpload        uint64
	totalDownload      uint64
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		prevNetCounters:    make(map[string]gnet.IOCountersStat),
		prevNetTime:        time.Now(),
		pingServers:        []string{"8.8.8.8", "1.1.1.1", "77.88.8.8"},
		currentServerIndex: 0,
	}
}

func (ms *MetricsService) GetSystemMetrics() models.SystemMetrics {
	return models.SystemMetrics{
		TimeStamp:      time.Now(),
		CPU:            ms.getCPUInfo(),
		Memory:         ms.getMemoryInfo(),
		Disk:           ms.getDiskInfo(),
		System:         ms.getSystemInfo(),
		Network:        ms.getNetworkMetrics(),
		Processes:      ms.getProcessMetrics(),
		NetworkDetails: ms.getNetworkDetailsMetrics(),
	}
}

func (ms *MetricsService) getCPUInfo() models.CPUInfo {

	cpuPercent, _ := cpu.Percent(400*time.Millisecond, false)
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

// ─── Network Stats ─────────────────────────────────────────────────────────

func (ms *MetricsService) calculateNetworkSpeed() (float64, float64, uint64, uint64) {
	currentCounters, err := gnet.IOCounters(true)
	if err != nil {
		return 0, 0, ms.totalUpload, ms.totalDownload
	}

	var curBytesSent, curBytesRecv uint64
	for _, counter := range currentCounters {
		curBytesRecv += counter.BytesRecv
		curBytesSent += counter.BytesSent
	}

	timeDelta := time.Since(ms.prevNetTime).Seconds()
	var uploadSpeed, downloadSpeed float64

	if len(ms.prevNetCounters) > 0 {
		var prevBytesRecv, prevBytesSent uint64
		for _, counter := range ms.prevNetCounters {
			prevBytesRecv += counter.BytesRecv
			prevBytesSent += counter.BytesSent
		}
		uploadSpeed = float64(curBytesSent-prevBytesSent) * 8 / timeDelta / 1000000   // Mb/s
		downloadSpeed = float64(curBytesRecv-prevBytesRecv) * 8 / timeDelta / 1000000 // Mb/s

		if curBytesSent > prevBytesSent {
			ms.totalUpload += (curBytesSent - prevBytesSent) / 1024 / 1024
			ms.totalDownload += (curBytesRecv - prevBytesRecv) / 1024 / 1024
		}
	}

	ms.prevNetCounters = make(map[string]gnet.IOCountersStat)
	for _, counter := range currentCounters {
		ms.prevNetCounters[counter.Name] = counter
	}
	ms.prevNetTime = time.Now()

	return uploadSpeed, downloadSpeed, ms.totalUpload, ms.totalDownload
}

func (ms *MetricsService) measurePing() float64 {
	server := ms.pingServers[ms.currentServerIndex]

	pingTime, err := ms.tcpPing(server)
	if err != nil {
		ms.currentServerIndex = (ms.currentServerIndex + 1) % len(ms.pingServers)
		return 1000.0
	}

	return pingTime
}

func (ms *MetricsService) tcpPing(server string) (float64, error) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", server+":80", 3*time.Second)
	if err != nil {
		return 0, fmt.Errorf("failed to connect ping server: %w", err)
	}
	defer conn.Close()

	ping := time.Since(start).Seconds() * 1000 // ping in ms

	return ping, nil
}

func (ms *MetricsService) getNetworkMetrics() models.NetworkStats {
	stats := models.NetworkStats{}

	upload, download, _, _ := ms.calculateNetworkSpeed()
	stats.CurrentUpload = upload
	stats.CurrentDownload = download

	stats.Ping = ms.measurePing()

	stats.IsOnline = stats.Ping <= 1000

	stats.LocalIP = ms.getLocalIP()

	return stats
}

func (ms *MetricsService) getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// ─── Network Details ─────────────────────────────────────────────────────────

func (ms *MetricsService) getNetworkDetailsMetrics() *models.NetworkDetails {
	_, _, totalUpload, totalDownload := ms.calculateNetworkSpeed()
	details := models.NetworkDetails{
		PublicIP:      getPublicIP(),
		MACAddress:    getMacAddr(),
		TotalUpload:   totalUpload,
		TotalDownload: totalDownload,
		ErrorMessage:  "",
	}

	return &details
}

func getMacAddr() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			return iface.HardwareAddr.String()
		}
	}

	return ""
}

func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Printf("failed to get public IP from API: %v", err)
		return ""
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read public IP from API: %v", err)
		return ""
	}

	return string(ip)
}

// ─── Process Metrics ─────────────────────────────────────────────────────────

func (ms *MetricsService) getProcessMetrics() []models.ProcessInfo {
	processes, err := process.Processes()
	if err != nil {
		return []models.ProcessInfo{}
	}

	var processMetrics []models.ProcessInfo
	for _, process := range processes {
		if info, err := ms.getSingleProcesInfo(process); err == nil {
			processMetrics = append(processMetrics, info)
		}
	}

	return processMetrics
}

func (ms *MetricsService) getSingleProcesInfo(p *process.Process) (models.ProcessInfo, error) {

	status, err := p.Status()
	if err != nil {
		return models.ProcessInfo{}, err
	}

	if status[0] == "Z" {
		return models.ProcessInfo{}, fmt.Errorf("zomvie process")
	}

	info := models.ProcessInfo{PID: p.Pid}

	process, err := p.Name()
	if err == nil {
		info.Process = process
	}
	cpuPercent, err := p.CPUPercent()
	if err == nil {
		info.CPUPercent = cpuPercent
	}
	memPercent, err := p.MemoryPercent()
	if err == nil {
		info.MemoryPercent = float64(memPercent)
	}
	memRSS, err := p.MemoryInfo()
	if err == nil {
		info.MemoryRSS = memRSS.RSS
	}
	info.Status = status[0]

	commandLine, err := p.Cmdline()
	if err == nil {
		if commandLine == "" {
			info.CommandLine = "System process"
		} else {
			info.CommandLine = commandLine
		}
	} else {
		info.CommandLine = "N/A"
	}

	user, err := p.Username()
	if err == nil {
		info.User = user
	}

	threads, err := p.NumThreads()
	if err == nil {
		info.Threads = int(threads)
	} else {
		info.Threads = -1
	}

	createTime, err := p.CreateTime()
	if err == nil {
		info.CreateTime = createTime / 1000
	} else {
		info.CreateTime = 0
	}

	return info, nil
}
