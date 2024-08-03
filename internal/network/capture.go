package network

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"homenethelper/pkg/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Capture struct {
	interfaceName string
	handle        *pcap.Handle
	StatsChan     chan map[string]*DeviceStats
	stopChan      chan struct{}
	isWifi        bool
	deviceCache   map[string]string
	cacheMutex    sync.RWMutex
}

func NewCapture(interfaceName string) (*Capture, error) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("error opening interface %s: %v", interfaceName, err)
	}

	isWifi := strings.HasPrefix(interfaceName, "wl") || strings.Contains(interfaceName, "wifi")

	return &Capture{
		interfaceName: interfaceName,
		handle:        handle,
		StatsChan:     make(chan map[string]*DeviceStats, 1), // Buffered channel
		stopChan:      make(chan struct{}),
		isWifi:        isWifi,
		deviceCache:   make(map[string]string),
	}, nil
}

func (c *Capture) Start() {
	go func() {
		packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
		stats := make(map[string]*DeviceStats)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case packet := <-packetSource.Packets():
				c.processPacket(packet, stats)
			case <-ticker.C:
				c.StatsChan <- stats
				stats = make(map[string]*DeviceStats) // Reset stats
			case <-c.stopChan:
				return
			}
		}
	}()
}

func (c *Capture) Stop() {
	close(c.stopChan)
	c.handle.Close()
}

func (c *Capture) processPacket(packet gopacket.Packet, stats map[string]*DeviceStats) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		log.Printf("Packet on %s has no network layer", c.interfaceName)
		return
	}

	srcIP := networkLayer.NetworkFlow().Src().String()
	dstIP := networkLayer.NetworkFlow().Dst().String()

	updateDeviceStats := func(ip string, isSent bool) {
		if _, ok := stats[ip]; !ok {
			deviceName := c.getDeviceName(ip)
			stats[ip] = &DeviceStats{
				DeviceName:     deviceName,
				TopContent:     make(map[string]*ContentStats),
				ConnectionType: c.getConnectionType(),
				Interface:      c.interfaceName,
			}
		}
		updateStats(stats[ip], packet, isSent)
	}

	if utils.IsPrivateIP(srcIP) {
		updateDeviceStats(srcIP, true)
	}

	if utils.IsPrivateIP(dstIP) {
		updateDeviceStats(dstIP, false)
	}
}

func (c *Capture) getDeviceName(ip string) string {
	c.cacheMutex.RLock()
	if name, ok := c.deviceCache[ip]; ok {
		c.cacheMutex.RUnlock()
		return name
	}
	c.cacheMutex.RUnlock()

	name := utils.LookupHostname(ip)
	
	c.cacheMutex.Lock()
	c.deviceCache[ip] = name
	c.cacheMutex.Unlock()

	return name
}

func (c *Capture) getConnectionType() string {
	if c.isWifi {
		return "Wi-Fi"
	}
	return "Ethernet"
}
