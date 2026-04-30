package pdf_helper

import (
	"fmt"
	"strings"
)

func amount(v float64) string {
	// make "1234567.89"
	s := fmt.Sprintf("%.2f", v)

	// split int/frac
	intPart, frac := s, "00"
	if i := strings.LastIndexByte(s, '.'); i >= 0 {
		intPart, frac = s[:i], s[i+1:]
	}

	// handle sign
	sign := ""
	if strings.HasPrefix(intPart, "-") {
		sign = "-"
		intPart = intPart[1:]
	}

	// insert thousands '.' into intPart (right-to-left)
	var out []byte
	cnt := 0
	for i := len(intPart) - 1; i >= 0; i-- {
		out = append(out, intPart[i])
		cnt++
		if cnt%3 == 0 && i != 0 {
			out = append(out, '.')
		}
	}
	// reverse out
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	// decimal separator ','
	return sign + string(out) + "," + frac
}
