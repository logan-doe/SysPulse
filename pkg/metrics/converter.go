package metrics

import "fmt"

func BytesToGB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024 / 1024
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d Bytes\n", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGPTE"[exp])
}
