package utils

import (
	"net"
)

// IsPrivateIP checks if the given IP address is a private address
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsPrivate()
}

// LookupHostname attempts to get the hostname for a given IP address
func LookupHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "Unknown"
	}
	return names[0]
}
