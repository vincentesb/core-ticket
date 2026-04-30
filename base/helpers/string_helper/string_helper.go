package string_helper

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// GetAbbreviation returns the abbreviation of the uppercase letter in the string.
// If no uppercase letter is found at the beginning of the string, it returns empty string.
func GetAbbreviation(str string) string {
	var result string
	for _, char := range str {
		if unicode.IsUpper(char) {
			result += string(char)
		}
	}
	return result
}

// Reverse used to return a reversed string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsUSKeyboardCharInput(s string) bool {
	re := regexp.MustCompile(`[^a-zA-Z0-9!"#$%&'()*+,\-.\/:;<=>?@\[\]\\^_` + "`" + `{|}~\t\s]`)
	tmpInput := re.ReplaceAllString(s, "")
	match, _ := regexp.MatchString(
		`^[a-zA-Z0-9!"#$%&'()*+,\-.\/:;<=>?@\[\]\\^_`+"`"+`{|}~\t\s]*$`,
		tmpInput,
	)
	return match
}

// Function to check if the string is alphanumeric and remove whitespaces first
func IsAlphanumericNoSpaces(s string) bool {
	s = strings.ReplaceAll(s, " ", "")

	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}

	return true
}

// UpperCaseFirst Uppercase the first character of the string
func UpperCaseFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func Unique(input []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string

	for _, str := range input {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = true
			result = append(result, str)
		}
	}

	return result
}

func addThousandsSeparator(numberStr, sep string) string {
	var formatted string
	for i, runeValue := range numberStr {
		if i > 0 && (len(numberStr)-i)%3 == 0 {
			formatted += sep
		}
		formatted += string(runeValue)
	}
	return formatted
}

func UcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

/*
SummarizeSliceByThreshold is a function for formatting a list of entities into a string.

	Example Output:


	if slice contains {"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"} and the threshold is 3, the output will be "New York, Los Angeles, Chicago and 2 more".
*/
func SummarizeSliceByThreshold(entities []string, threshold int) string {
	if len(entities) == 0 {
		return ""
	}
	if len(entities) <= threshold {
		return strings.Join(entities, ", ")
	} else {
		return fmt.Sprintf("%s and %d more", strings.Join(entities[:threshold], ", "), len(entities)-threshold)
	}
}

/*
	AddSpaceOnUCWords returns the characters being spaced on each uppercase letter in the string.

# If the character is uppercase, and it's not the first character, adds a space before it

example usage:
camelCase => camel case
*/
func AddSpaceOnUCWords(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 {
			result = append(result, ' ')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func ConvertCamelCaseToNormal(input string) string {
	var words []string
	var currentWord []rune

	for _, r := range input {
		if unicode.IsUpper(r) && len(currentWord) > 0 {
			// Append the current word to words slice
			words = append(words, string(currentWord))
			// Start a new word
			currentWord = []rune{r}
		} else {
			currentWord = append(currentWord, r)
		}
	}

	// Append the last word
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}

	// Capitalize the first letter of each word and join them with spaces
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	return strings.Join(words, " ")
}

/**
 * ToArray converts a string to an array of strings.
 * It splits the string by the specified separator and trims each part.
 * If the separator is empty, it splits the string by commas.
 * e.g. "a,,c" becomes ["a", "c"]
 */
func ToArray(s string, separator string) []string {
	if separator == "" {
		separator = ","
	}

	parts := strings.Split(s, separator)

	var cleanedParts []string

	for _, part := range parts {
		if part == "" {
			continue
		}

		cleanedParts = append(cleanedParts, strings.TrimSpace(part))
	}

	return cleanedParts
}

// NormalizeNumberString normalizes a number string by handling both
// dot-as-thousands (e.g. "100.123,1") and comma-as-thousands (e.g. "100,123.1") conventions.
func NormalizeNumberString(s string) string {
	lastDot := strings.LastIndex(s, ".")
	lastComma := strings.LastIndex(s, ",")

	if lastDot > lastComma {
		// comma is thousands separator, e.g. "100,123.1"
		return strings.ReplaceAll(s, ",", "")
	} else if lastComma > lastDot {
		// dot is thousands separator, e.g. "100.123,1"
		s = strings.ReplaceAll(s, ".", "")
		return strings.ReplaceAll(s, ",", ".")
	}
	return s
}
