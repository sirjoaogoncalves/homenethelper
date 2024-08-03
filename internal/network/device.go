package network

import (
	"time"
)

type ContentStats struct {
	Protocol string
	Port     uint16
	Bytes    uint64
}

type DeviceStats struct {
	BytesReceived    uint64
	BytesSent        uint64
	LastUpdateTime   time.Time
	DownloadSpeed    float64 // in bytes per second
	UploadSpeed      float64 // in bytes per second
	TotalDownloaded  uint64
	TotalUploaded    uint64
	DeviceName       string
	TopContent       map[string]*ContentStats // Key: "Protocol:Port"
	ConnectionType   string                   // "Wi-Fi" or "Ethernet"
	Interface        string                   // The network interface this device was seen on
}
