package network

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func updateStats(stats *DeviceStats, packet gopacket.Packet, isSent bool) {
	bytes := uint64(len(packet.Data()))
	now := time.Now()

	duration := now.Sub(stats.LastUpdateTime).Seconds()
	if duration > 0 {
		if isSent {
			stats.UploadSpeed = float64(bytes) / duration
			stats.BytesSent += bytes
			stats.TotalUploaded += bytes
		} else {
			stats.DownloadSpeed = float64(bytes) / duration
			stats.BytesReceived += bytes
			stats.TotalDownloaded += bytes
		}
	}
	stats.LastUpdateTime = now

	updateContentStats(stats, packet)
}

func updateContentStats(stats *DeviceStats, packet gopacket.Packet) {
	protocol, port := getProtocolAndPort(packet)
	key := fmt.Sprintf("%s:%d", protocol, port)
	
	if _, ok := stats.TopContent[key]; !ok {
		stats.TopContent[key] = &ContentStats{
			Protocol: protocol,
			Port:     port,
		}
	}
	stats.TopContent[key].Bytes += uint64(len(packet.Data()))
}

func getProtocolAndPort(packet gopacket.Packet) (string, uint16) {
	// Start with default values
	protocol := "Unknown"
	var port uint16 = 0

	// Check for TCP layer
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		protocol = "TCP"
		tcp, _ := tcpLayer.(*layers.TCP)
		port = uint16(tcp.DstPort)
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		// Check for UDP layer
		protocol = "UDP"
		udp, _ := udpLayer.(*layers.UDP)
		port = uint16(udp.DstPort)
	} else if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		// Check for ICMP layer
		protocol = "ICMP"
	} else {
		log.Printf("Unknown protocol in packet: %v", packet)
	}

	return protocol, port
}

func MergeStats(dest, src map[string]*DeviceStats) {
	for ip, srcStats := range src {
		if destStats, exists := dest[ip]; exists {
			destStats.BytesReceived += srcStats.BytesReceived
			destStats.BytesSent += srcStats.BytesSent
			destStats.TotalDownloaded += srcStats.TotalDownloaded
			destStats.TotalUploaded += srcStats.TotalUploaded
			if srcStats.LastUpdateTime.After(destStats.LastUpdateTime) {
				destStats.LastUpdateTime = srcStats.LastUpdateTime
			}
			destStats.DownloadSpeed += srcStats.DownloadSpeed
			destStats.UploadSpeed += srcStats.UploadSpeed

			for key, srcContent := range srcStats.TopContent {
				if destContent, contentExists := destStats.TopContent[key]; contentExists {
					destContent.Bytes += srcContent.Bytes
				} else {
					destStats.TopContent[key] = &ContentStats{
						Protocol: srcContent.Protocol,
						Port:     srcContent.Port,
						Bytes:    srcContent.Bytes,
					}
				}
			}
		} else {
			dest[ip] = &DeviceStats{
				BytesReceived:    srcStats.BytesReceived,
				BytesSent:        srcStats.BytesSent,
				LastUpdateTime:   srcStats.LastUpdateTime,
				DownloadSpeed:    srcStats.DownloadSpeed,
				UploadSpeed:      srcStats.UploadSpeed,
				TotalDownloaded:  srcStats.TotalDownloaded,
				TotalUploaded:    srcStats.TotalUploaded,
				DeviceName:       srcStats.DeviceName,
				TopContent:       make(map[string]*ContentStats),
			}
			for key, content := range srcStats.TopContent {
				dest[ip].TopContent[key] = &ContentStats{
					Protocol: content.Protocol,
					Port:     content.Port,
					Bytes:    content.Bytes,
				}
			}
		}
	}
}
