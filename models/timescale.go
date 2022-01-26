package models

import "time"

// timescaleDB 超表定义

type TsDeviceLoad struct {
	Timestamp   time.Time `gorm:"primaryKey" json:"timestamp"`
	Sn          string    `gorm:"primaryKey" json:"sn"`
	Name        string    `json:"name"`
	CpuLoad     int64     `json:"cpu_load"`
	MemUsage    int64     `json:"mem_usage"`
	MemPercent  float64   `json:"mem_percent"`
	DiskUsage   int64     `json:"disk_usage"`
	DiskPercent float64   `json:"disk_percent"`
}
