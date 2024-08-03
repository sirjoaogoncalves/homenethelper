package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"homenethelper/internal/network"
	"homenethelper/internal/ui"
)

func main() {
	// Set up logging
	logFile, err := os.OpenFile("homenethelper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Parse command-line flags
	interfacesFlag := flag.String("interfaces", "", "Comma-separated list of network interfaces to monitor")
	refreshRateFlag := flag.Int("refresh", 5, "Refresh rate in seconds")
	flag.Parse()

	if *interfacesFlag == "" {
		fmt.Println("Please specify network interfaces to monitor.")
		fmt.Println("Usage: homenethelper -interfaces=eth0,wlan0,wlan1 [-refresh=5]")
		os.Exit(1)
	}

	interfaces := strings.Split(*interfacesFlag, ",")

	// Initialize packet capture for each interface
	captures := make([]*network.Capture, 0, len(interfaces))
	for _, iface := range interfaces {
		capture, err := network.NewCapture(iface)
		if err != nil {
			log.Printf("Failed to initialize packet capture on %s: %v", iface, err)
			continue
		}
		captures = append(captures, capture)
		log.Printf("Successfully initialized capture on interface: %s", iface)
	}

	if len(captures) == 0 {
		log.Fatal("No valid interfaces to capture")
	}

	// Initialize combined stats
	combinedStats := make(map[string]*network.DeviceStats)

	// Start UI
	ui := ui.NewTerminalUI()

	// Setup signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start capturing packets on all interfaces
	for _, capture := range captures {
		go capture.Start()
	}

	// Ticker for regular updates
	ticker := time.NewTicker(time.Duration(*refreshRateFlag) * time.Second)
	defer ticker.Stop()

	// Main loop
	for {
		select {
		case <-signalChan:
			fmt.Println("\nShutting down...")
			for _, capture := range captures {
				capture.Stop()
			}
			return
		case <-ticker.C:
			// Collect stats from all captures
			for _, capture := range captures {
				select {
				case stats := <-capture.StatsChan:
					network.MergeStats(combinedStats, stats)
				default:
					// No stats available, continue
				}
			}
			ui.Update(combinedStats)
		}
	}
}
