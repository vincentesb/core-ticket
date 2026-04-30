package file_helper

import (
	"path/filepath"
	"slices"
	"strings"
)

// ValidateFileExtension checks whether a given filename has a valid extension.
// It returns true if the extension is valid (i.e., not empty and found in the allowedExts list),
// and false otherwise.
//
// Parameters:
//   - filename: the full name of the file (e.g., "report.xlsx").
//   - allowedExts: a list of allowed lowercase file extensions without the leading dot (e.g., []string{"xlsx", "csv"}).
//
// Example:
//
//	valid := ValidateFileExtension("report.xlsx", []string{"xlsx", "csv"})
//	// valid == true
//
//	valid := ValidateFileExtension("report", []string{"xlsx", "csv"})
//	// valid == false
//
//	valid := ValidateFileExtension("report.exe", []string{"xlsx", "csv"})
//	// valid == false
func ValidateFileExtension(filename string, allowedExts []string) bool {
	ext := filepath.Ext(filename)
	return len(ext) >= 2 && slices.Contains(allowedExts, strings.ToLower(ext[1:]))
}
