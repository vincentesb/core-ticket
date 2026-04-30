package query_sanitizer

import (
	"regexp"
	"strings"
)

// List of allowed SQL functions for ORDER BY
var safeFunctions = []string{
	"FIELD", "CONCAT", "COALESCE", "ABS", "IF", "COUNT",
	"MAX", "MIN", "SUM", "AVG",
	"LENGTH", "LOWER", "UPPER", "TRIM",
	"LEFT", "RIGHT",
	"DATE", "YEAR", "MONTH", "DAY",
}

// SanitizeOrderBy cleans an ORDER BY clause to mitigate SQL injection risks.
// If dangerous patterns are detected, it returns a safe fallback like "1".
func SanitizeOrderBy(input string) string {
	input = sanitizeCommentsPrecise(input)
	input = removeDangerousKeywords(input)
	input = removeSubqueries(input)
	input = strings.ReplaceAll(input, ";", "")
	input = strings.TrimSpace(input)

	if isDangerous(input) {
		return "1" // fallback safe ORDER BY expression
	}

	return input
}

// sanitizeCommentsPrecise removes SQL comments that are outside of quotes.
func sanitizeCommentsPrecise(input string) string {
	var (
		inSingleQuote = false
		inDoubleQuote = false
		escaped       = false
		result        strings.Builder
		i             = 0
	)

	for i < len(input) {
		c := input[i]

		if c == '\\' && !escaped {
			escaped = true
			result.WriteByte(c)
			i++
			continue
		}

		if c == '\'' && !inDoubleQuote && !escaped {
			inSingleQuote = !inSingleQuote
		} else if c == '"' && !inSingleQuote && !escaped {
			inDoubleQuote = !inDoubleQuote
		}

		if !inSingleQuote && !inDoubleQuote && !escaped {
			if c == '-' && i+1 < len(input) && input[i+1] == '-' {
				break
			}
			if c == '#' {
				break
			}
		}

		result.WriteByte(c)
		escaped = false
		i++
	}

	return strings.TrimSpace(result.String())
}

// removeDangerousKeywords strips risky SQL keywords
func removeDangerousKeywords(input string) string {
	pattern := regexp.MustCompile(`(?i)\b(DROP|DELETE|INSERT|UPDATE|UNION|ALTER|TRUNCATE|CREATE|REPLACE|EXEC|EXECUTE|SLEEP)\b`)
	return pattern.ReplaceAllString(input, "")
}

// removeSubqueries strips simple subqueries like (SELECT ...)
func removeSubqueries(input string) string {
	subQueryPattern := regexp.MustCompile(`(?i)\(SELECT[^()]*\)`)
	return subQueryPattern.ReplaceAllString(input, "")
}

// isDangerous checks if input contains anything outside the expected pattern
func isDangerous(input string) bool {
	safeFuncSet := make(map[string]bool)
	for _, fn := range safeFunctions {
		safeFuncSet[strings.ToUpper(fn)] = true
	}

	// Check for used functions in the input
	funcPattern := regexp.MustCompile(`(?i)([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	matches := funcPattern.FindAllStringSubmatch(input, -1)

	for _, m := range matches {
		funcName := strings.ToUpper(m[1])
		if !safeFuncSet[funcName] {
			return true // unknown or unsafe function used
		}
	}

	return false
}
