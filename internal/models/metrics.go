package models

import "time"

type SystemMetrics struct {
	TimeStamp time.Time  `json:"timestamp"`
	CPU       CPUInfo    `json:"cpu"`
	Memory    MemInfo    `json:"memory"`
	Disk      DiskInfo   `json:"disk"`
	System    SystemInfo `json:"system"`
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

type Alert struct { // alert is to show some important system messages to user
	ID        string    `json:"id"`        // id
	Type      string    `json:"type"`      // cpu, ram or disk
	Message   string    `json:"message"`   // alert text
	Level     string    `json:"level"`     // alert level (warning, critical, etc.)
	Value     float64   `json:"value"`     // current value
	Threshold float64   `json:"threshold"` // use to determine level of alert
	Timestamp time.Time `json:"timestamp"` // time LOL
}
