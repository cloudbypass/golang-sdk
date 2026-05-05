package cloudbypass

import (
	"fmt"
	"os"
	"strings"
)

func getEnv(key string, default_value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return default_value
}

var convertUnits = []string{"", "K", "M", "G", "T", "P", "E", "Z", "Y"}

// ConvertBytes formats a byte count as a human-readable string (1024-based).
// endUnit is one of "", "K", "M", …, "Y"; empty defaults to "Y". Single runes are uppercased.
func ConvertBytes(value float64, endUnit string) (string, error) {
	if value == 0 {
		return "0", nil
	}
	if endUnit == "" {
		endUnit = "Y"
	}
	if len(endUnit) == 1 {
		endUnit = strings.ToUpper(endUnit)
	}
	endIdx := -1
	for i, u := range convertUnits {
		if u == endUnit {
			endIdx = i
			break
		}
	}
	if endIdx < 0 {
		return "", fmt.Errorf("invalid endUnit: %s", endUnit)
	}
	unit := 0
	v := value
	for v >= 1024 {
		v /= 1024
		unit++
		if unit == endIdx {
			break
		}
	}
	return fmt.Sprintf("%.2f %sB", v, convertUnits[unit]), nil
}
