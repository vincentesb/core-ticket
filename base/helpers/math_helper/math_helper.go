package math_helper

import "math"

func Round(val float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(val*shift) / shift
}

// Round4 performs PHP-like "round half up" rounding to 4 decimal places. Equivalent to number_format($hpp, 4, '.', ”)
func Round4(v float64) float64 {
	factor := math.Pow(10, 4)
	return math.Floor(v*factor+0.5) / factor
}
