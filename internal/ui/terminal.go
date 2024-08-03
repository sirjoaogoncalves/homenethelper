package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"homenethelper/internal/network"
	"homenethelper/pkg/utils"
)

type TerminalUI struct {
	lastUpdate  time.Time
	refreshRate time.Duration
}

func NewTerminalUI() *TerminalUI {
	return &TerminalUI{
		refreshRate: 5 * time.Second, // Update every 5 seconds
	}
}

func (t *TerminalUI) Update(stats map[string]*network.DeviceStats) {
	now := time.Now()
	if now.Sub(t.lastUpdate) < t.refreshRate {
		return
	}
	t.lastUpdate = now

	// Clear the screen and move cursor to top-left
	fmt.Print("\033[2J\033[H")

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "IP Address\tDevice Name\tConnection\tInterface\tDown Speed\tUp Speed\tTotal Down\tTotal Up\tTop Protocol\tTop Port")
	fmt.Fprintln(w, strings.Repeat("-", 150)) // Separator line

	// Sort IPs for consistent output
	var ips []string
	for ip := range stats {
		ips = append(ips, ip)
	}
	sort.Strings(ips)

	for _, ip := range ips {
		stat := stats[ip]
		topProtocol, topPort := getTopContent(stat.TopContent)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s/s\t%s/s\t%s\t%s\t%s\t%s\n",
			ip,
			stat.DeviceName,
			stat.ConnectionType,
			stat.Interface,
			utils.BytesToHumanReadable(uint64(stat.DownloadSpeed)),
			utils.BytesToHumanReadable(uint64(stat.UploadSpeed)),
			utils.BytesToHumanReadable(stat.TotalDownloaded),
			utils.BytesToHumanReadable(stat.TotalUploaded),
			topProtocol,
			topPort)
	}

	w.Flush()

	// Print last update time
	fmt.Printf("\nLast updated: %s\n", now.Format("15:04:05"))
	fmt.Println("Press Ctrl+C to exit")
}

func getTopContent(content map[string]*network.ContentStats) (string, string) {
	if len(content) == 0 {
		return "N/A", "N/A"
	}
	var topProtocol string
	var topPort uint16
	var topBytes uint64
	for _, stats := range content {
		if stats.Bytes > topBytes {
			topProtocol = stats.Protocol
			topPort = stats.Port
			topBytes = stats.Bytes
		}
	}
	return fmt.Sprintf("%s (%s)", topProtocol, utils.BytesToHumanReadable(topBytes)), 
	       fmt.Sprintf("%d (%s)", topPort, utils.BytesToHumanReadable(topBytes))
}
