package decimal_helper

import (
	"fmt"
	"math"
	"strings"
)

// IsDecimal Check if the float64 value contains a decimal part
func IsDecimal(f float64) bool {
	return f != math.Trunc(f)
}

// FormatDecimal used for decimal format to string (thousand separator)
//
// example value: 10900,80
//
// example response : 10,900.80
func FormatDecimal(value float64) string {
	// Format the number with two decimal places
	formatted := fmt.Sprintf("%.2f", value)

	// Insert thousands separator
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	// Add commas to the integer part
	for i := len(integerPart) - 3; i > 0; i -= 3 {
		integerPart = integerPart[:i] + "," + integerPart[i:]
	}

	return integerPart + "." + decimalPart
}

// RoundFloat64 rounds a float64 number x to a given number of decimal places.
// Example: RoundFloat64(3.14159265, 4) returns 3.1416.
func RoundFloat64(x float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(x*pow) / pow
}

func CountDecimals(f float64) int {
	s := fmt.Sprintf("%.10f", f)

	// trim trailing zeros.
	s = strings.TrimRight(s, "0")

	// Split by the decimal point
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return 0
	}

	return len(parts[1])
}
