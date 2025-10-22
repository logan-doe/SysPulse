package models

import (
	"time"
)

type SystemMetrics struct {
	TimeStamp      time.Time       `json:"timestamp"`
	CPU            CPUInfo         `json:"cpu"`
	Memory         MemInfo         `json:"memory"`
	Disk           DiskInfo        `json:"disk"`
	System         SystemInfo      `json:"system"`
	Alerts         []Alert         `json:"alerts,omitempty"`
	Network        NetworkStats    `json:"network"`
	Processes      []ProcessInfo   `json:"processes"`
	NetworkDetails *NetworkDetails `json:"network_details,omitempty"`
}

type CPUInfo struct {
	Usage  float64 `json:"usage"`  // cpu usage in %
	Cores  int     `json:"cores"`  // amount of cores
	Load1  float64 `json:"load1"`  // load avg for 1 min
	Load5  float64 `json:"load5"`  // load avg for 5 mins
	Load15 float64 `json:"load15"` // load avg for 15 mins
}

type MemInfo struct {
	Total     uint64  `json:"total"`     // total memory size in bytes
	Used      uint64  `json:"used"`      // used memory in bytes
	Available uint64  `json:"available"` // available memory in bytes
	Usage     float64 `json:"usage"`     // memory usage in %
}

type DiskInfo struct {
	Total uint64  `json:"total"` // total disk volume in bytes
	Used  uint64  `json:"used"`  // used space in bytes
	Free  uint64  `json:"free"`  // free spaces in bytes
	Usage float64 `json:"usage"` // disk usage in %
}

type SystemInfo struct {
	Hostname string `json:"hostname"` // name of userspace
	OS       string `json:"os"`       // os name
	Platform string `json:"platform"` // platform name
	Uptime   uint64 `json:"uptime"`   // time since system started
}

// system alerts
type Alert struct {
	ID        string    `json:"id"`        // id
	Type      string    `json:"type"`      // cpu, ram or disk
	Message   string    `json:"message"`   // alert text
	Level     string    `json:"level"`     // alert level (warning, critical, etc.)
	Value     float64   `json:"value"`     // current value
	Threshold float64   `json:"threshold"` // use to determine level of alert
	Timestamp time.Time `json:"timestamp"` // time LOL
	Active    bool      `json:"active"`    // show if alertis active
}

// alert config
type AlertConfig struct {
	CPUTreshold  float64 `json:"cpu_treshold"`
	RAMTreshold  float64 `json:"ram_treshold"`
	DiskTreshold float64 `json:"disk_treshold"`
	Enabled      bool    `json:"enabled"`
}

// alert hystory
type AlertHistory struct {
	Alerts []Alert `json:"alerts"`
	Stats  struct {
		TotalAlerts  int `json:"total_alerts"`
		ActiveAlerts int `json:"active_alerts"`
		TodayAlerts  int `json:"today_alerts"`
	} `json:"stats"`
}

type ProcessInfo struct {
	PID           int32   `json:"pid"`
	Process       string  `json:"process"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryRSS     uint64  `json:"memory_rss"`
	Status        string  `json:"status"`
	CommandLine   string  `json:"commandline"`
	User          string  `json:"user"`
	CreateTime    int64   `json:"createtime"`
	Threads       int     `json:"threads"`
}

type NetworkStats struct {
	IsOnline        bool    `json:"is_online"`        // am i online ?
	CurrentUpload   float64 `json:"current_upload"`   // net speed in Mb/s
	CurrentDownload float64 `json:"current_download"` // net speed in Mb/s
	Ping            float64 `json:"ping"`             // net dealy in milliseconds
	LocalIP         string  `json:"local_ip"`
}

type NetworkDetails struct {
	PublicIP      string `json:"public_ip"`
	MACAddress    string `json:"mac_address"`
	TotalUpload   uint64 `json:"total_upload"`            // total value of used traffic in Mb
	TotalDownload uint64 `json:"total_download"`          // total value of used traffic in Mb
	ErrorMessage  string `json:"error_message,omitempty"` // any error ?
}

type PingStrategy struct {
	PrinaryServers  []string // reliable servers
	FallbackServers []string // additional servers if any problem with primaries
	TimeoutMS       int
}
