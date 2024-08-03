package utils

import (
	"fmt"
)

// BytesToHumanReadable converts bytes to a human-readable string
func BytesToHumanReadable(bytes uint64) string {
	return bytesToHumanReadableHelper(float64(bytes))
}

func bytesToHumanReadableHelper(bytes float64) string {
	const unit = 1024.0
	if bytes < unit {
		return fmt.Sprintf("%.1f B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", bytes/div, "KMGTPE"[exp])
}
