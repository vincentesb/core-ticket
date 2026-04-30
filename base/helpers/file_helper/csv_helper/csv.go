package csv_helper

import (
	"bytes"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/constants/error_code"
	"encoding/csv"
	"fmt"
	"os"
)

const SeparatorCsv = ";"

func SaveCSVInTemp(fileName string, data [][]string, separator rune) (string, error) {
	f, err := os.Create(fmt.Sprintf(
		"%s%s",
		os.TempDir(),
		fileName,
	))
	if err != nil {
		return "", err
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	w.Comma = separator

	for _, d := range data {
		if err := w.Write(d); err != nil {
			return "", err
		}
	}

	return f.Name(), nil
}

// WriteToCsv used to write data into csv buffer. It returns buffer and error
func WriteToCsv(data [][]string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	writer := csv.NewWriter(&buf)
	// Write data to the CSV writer
	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return buf, error_helper.New(fmt.Errorf("error writing record to CSV: %v", err), error_code.UnknownError)
		}
	}

	// Flush the writer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return buf, error_helper.New(fmt.Errorf("error flushing CSV writer: %v", err), error_code.UnknownError)
	}

	return buf, err
}
