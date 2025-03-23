package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// FormatIDR formats a number to Indonesian Rupiah format
func FormatIDR(amount float64) string {
	// Handle negative amounts
	isNegative := amount < 0
	if isNegative {
		amount = math.Abs(amount)
	}

	// Convert to string and split decimal part if exists
	amountStr := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(amountStr, ".")
	wholePart := parts[0]
	decimalPart := "00"
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Add thousand separators
	length := len(wholePart)
	var result strings.Builder

	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteByte(wholePart[i])
	}

	// Combine all parts
	formatted := fmt.Sprintf("Rp %s,%s", result.String(), decimalPart)
	if isNegative {
		return "-" + formatted
	}
	return formatted
}

// ParseIDR parses an IDR formatted string back to float64
func ParseIDR(idr string) (float64, error) {
	// Remove currency symbol, thousand separators and spaces
	idr = strings.TrimSpace(idr)
	idr = strings.TrimPrefix(idr, "Rp")
	idr = strings.TrimPrefix(idr, "Rp.")
	idr = strings.TrimSpace(idr)
	idr = strings.ReplaceAll(idr, ".", "")
	idr = strings.ReplaceAll(idr, ",", ".")

	// Parse to float64
	return strconv.ParseFloat(idr, 64)
}

// FormatPercentage formats a decimal number to percentage string
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value*100)
}

// ParsePercentage parses a percentage string back to float64
func ParsePercentage(percentage string) (float64, error) {
	// Remove % symbol and spaces
	percentage = strings.TrimSpace(percentage)
	percentage = strings.TrimSuffix(percentage, "%")
	percentage = strings.TrimSpace(percentage)

	// Parse to float64 and convert to decimal
	value, err := strconv.ParseFloat(percentage, 64)
	if err != nil {
		return 0, err
	}
	return value / 100, nil
}
