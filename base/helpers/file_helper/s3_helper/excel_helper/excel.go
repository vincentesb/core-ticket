package excel_helper

import (
	"bytes"
	"core-ticket/base/helpers/error_helper"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

/*
- ErrorValidator: error from validator
- SheetName: sheet name from excel
- LastColumnName: filled with last column of excel (used to set column of error message). Example: "L" -> column L in excel
- StartIndex: start index row of excel (example: excel header is 1, then start index will be 0)
*/
type ErrorFileSheetMeta struct {
	ErrorValidator *error_helper.Error
	SheetName      string
	LastColumnName string
	StartIndex     int
}

// ExtractExcelFileFromPath : startIndex is header row number, start from 0
func ExtractExcelFileFromPath[T interface{}](filePath string, startIndex int, opts ...excelize.Options) ([]T, error) {
	stream, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(stream *os.File) {
		_ = stream.Close()
	}(stream)

	xlsx, errXlsx := excelize.OpenReader(stream, opts...)
	if errXlsx != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", errXlsx)
	}

	return extractExcelRows[T](xlsx, "", startIndex)
}

func ExtractExcelFileSpecifiedSheetFromPath[T interface{}](filePath string, sheetName string, startIndex int, opts ...excelize.Options) ([]T, error) {
	fileContent, fileErr := os.Open(filePath)
	if fileErr != nil {
		return nil, fmt.Errorf("failed to open file: %v", fileErr)
	}
	defer func(fileContent multipart.File) {
		_ = fileContent.Close()
	}(fileContent)

	fileBytes, readErr := io.ReadAll(fileContent)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read file content: %v", readErr)
	}

	xlsx, errXlsx := excelize.OpenReader(strings.NewReader(string(fileBytes)), opts...)
	if errXlsx != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", errXlsx)
	}

	return extractExcelRows[T](xlsx, sheetName, startIndex) // Extract from specific sheet
}

// ExtractExcelFile : startIndex is header row number, start from 0
func ExtractExcelFile[T interface{}](file *multipart.FileHeader, startIndex int, opts ...excelize.Options) ([]T, error) {
	fileContent, fileErr := file.Open()
	if fileErr != nil {
		return nil, fmt.Errorf("failed to open file: %v", fileErr)
	}
	defer func(fileContent multipart.File) {
		_ = fileContent.Close()
	}(fileContent)

	fileBytes, readErr := io.ReadAll(fileContent)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read file content: %v", readErr)
	}

	xlsx, errXlsx := excelize.OpenReader(strings.NewReader(string(fileBytes)), opts...)
	if errXlsx != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", errXlsx)
	}

	return extractExcelRows[T](xlsx, "", startIndex)
}

// ExtractExcelFileSpecifiedSheet : startIndex is header row number, start from 0
// used to extract excel file based on sheetName in params
func ExtractExcelFileSpecifiedSheet[T interface{}](file *multipart.FileHeader, sheetName string, startIndex int, opts ...excelize.Options) ([]T, error) {
	fileContent, fileErr := file.Open()
	if fileErr != nil {
		return nil, fmt.Errorf("failed to open file: %v", fileErr)
	}
	defer func(fileContent multipart.File) {
		_ = fileContent.Close()
	}(fileContent)

	fileBytes, readErr := io.ReadAll(fileContent)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read file content: %v", readErr)
	}

	xlsx, errXlsx := excelize.OpenReader(strings.NewReader(string(fileBytes)), opts...)
	if errXlsx != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", errXlsx)
	}

	return extractExcelRows[T](xlsx, sheetName, startIndex) // Extract from specific sheet
}

// ExtractExcelRowsBySheetIdx
// - used to extract Excel rows based on sheetIndex in params.
// - startIndex is header row number, start from 0
func ExtractExcelRowsBySheetIdx[T interface{}](file *multipart.FileHeader, sheetIndex int, startIndex int, opts ...excelize.Options) ([]T, error) {
	fileContent, fileErr := file.Open()
	if fileErr != nil {
		return nil, fmt.Errorf("failed to open file: %v", fileErr)
	}
	defer func(fileContent multipart.File) {
		_ = fileContent.Close()
	}(fileContent)

	fileBytes, readErr := io.ReadAll(fileContent)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read file content: %v", readErr)
	}

	xlsx, errXlsx := excelize.OpenReader(strings.NewReader(string(fileBytes)), opts...)
	if errXlsx != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", errXlsx)
	}

	sheetName := xlsx.GetSheetName(sheetIndex)

	return extractExcelRows[T](xlsx, sheetName, startIndex) // Extract from specific sheet
}

// extractExcelRows : Extracts Excel rows from a specific sheet or all sheets starting from the given startIndex
func extractExcelRows[T interface{}](xlsx *excelize.File, sheetName string, startIndex int) ([]T, error) {
	var rowDataList []T

	if sheetName != "" {
		// Extract rows from specific sheet
		rows, err := xlsx.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("error getting rows from sheet %s: %v", sheetName, err)
		}
		rows = processEmptyRows(rows)
		rowDataList, err = processRows[T](rows, startIndex)
		if err != nil {
			return nil, err
		}
	} else {
		// Extract rows from all sheets
		for _, sheet := range xlsx.GetSheetList() {
			rows, err := xlsx.GetRows(sheet)
			if err != nil {
				return nil, fmt.Errorf("error getting rows from sheet %s: %v", sheet, err)
			}
			rows = processEmptyRows(rows)
			tempDataList, err := processRows[T](rows, startIndex)
			if err != nil {
				return nil, err
			}
			rowDataList = append(rowDataList, tempDataList...)
		}
	}

	return rowDataList, nil
}

// processEmptyRows : Removes empty rows from the provided slice of strings
func processEmptyRows(rows [][]string) [][]string {
	var filteredRows [][]string
	for _, row := range rows {
		isEmpty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			filteredRows = append(filteredRows, row)
		}
	}
	return filteredRows
}

// processRows : Processes rows into the provided struct type T
func processRows[T interface{}](rows [][]string, startIndex int) ([]T, error) {
	var rowDataList []T

	for i, row := range rows {
		if i > startIndex {
			var temp T
			rowData := reflect.ValueOf(&temp).Elem()
			for j, cell := range row {
				if j < rowData.NumField() {
					field := rowData.Field(j)
					if field.IsValid() {
						switch field.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							v, _ := strconv.ParseInt(cell, 10, 64)
							field.SetInt(v)
						case reflect.Float32, reflect.Float64:
							v, _ := strconv.ParseFloat(cell, 64)
							field.SetFloat(v)
						case reflect.Bool:
							v, _ := strconv.ParseBool(cell)
							field.SetBool(v)
						case reflect.String:
							field.SetString(cell)
						default:
							panic("unhandled default case for excel")
						}
					}
				} else {
					break
				}
			}
			rowDataList = append(rowDataList, temp)
		}
	}

	return rowDataList, nil
}

// WriteXLSXUploadByPath : startIndex is header row number, start from 0
func WriteXLSXUploadByPath(filePath string, startIndex int, lastColumnName string, errs []*error_helper.ValidationError) (bytes.Buffer, error) {
	stream, err := os.Open(filePath)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer stream.Close()

	xlsx, err := excelize.OpenReader(stream)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// Write error messages to the Excel file
	return writeErrorMessagesToExcel(xlsx, startIndex, lastColumnName, errs)
}

// WriteXLSXUpload : startIndex is header row number, start from 0
func WriteXLSXUpload(file *multipart.FileHeader, startIndex int, lastColumnName string, errs []*error_helper.ValidationError) (bytes.Buffer, error) {
	stream, err := file.Open()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer stream.Close()

	xlsx, err := excelize.OpenReader(stream)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// Write error messages to the Excel file
	return writeErrorMessagesToExcel(xlsx, startIndex, lastColumnName, errs)
}

// WriteXLSXUploadMultipleSheet used to write error message in excel with multiple sheets
func WriteXLSXUploadMultipleSheet(file *multipart.FileHeader, sheetMetas []ErrorFileSheetMeta) (bytes.Buffer, error) {
	stream, err := file.Open()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer stream.Close()

	xlsx, err := excelize.OpenReader(stream)
	if err != nil {
		return bytes.Buffer{}, err
	}

	var buf bytes.Buffer
	for _, sheetMeta := range sheetMetas {
		groupedMessages := groupValidationErrors(sheetMeta.ErrorValidator.ValidationErrors())

		// Write grouped messages to the specific sheet
		for row, message := range groupedMessages {
			_ = xlsx.SetCellValue(sheetMeta.SheetName, fmt.Sprintf("%s%d", sheetMeta.LastColumnName, row+sheetMeta.StartIndex+1), message)
		}
	}

	if err := xlsx.Write(&buf); err != nil {
		return buf, err
	}

	return buf, nil
}

// Helper function to group validation errors by row
func groupValidationErrors(errs []*error_helper.ValidationError) map[int]string {
	groupedMessages := make(map[int]string)
	for _, vErr := range errs {
		if msg, exists := groupedMessages[vErr.Row]; exists {
			// Append the new message to the existing value
			groupedMessages[vErr.Row] = msg + ", " + vErr.Error()
		} else {
			// Add the message to the map
			groupedMessages[vErr.Row] = vErr.Error()
		}
	}
	return groupedMessages
}

// Helper function to write error messages to Excel
func writeErrorMessagesToExcel(xlsx *excelize.File, startIndex int, lastColumnName string, errs []*error_helper.ValidationError) (bytes.Buffer, error) {
	groupedMessages := groupValidationErrors(errs)

	// Write grouped messages to the first sheet
	for row, message := range groupedMessages {
		_ = xlsx.SetCellValue(xlsx.GetSheetList()[0], fmt.Sprintf("%s%d", lastColumnName, row+startIndex+1), message)
	}

	var buf bytes.Buffer
	if err := xlsx.Write(&buf); err != nil {
		return buf, err
	}

	return buf, nil
}

func InsertHeader[T interface{}](xlsx *excelize.File, row int, sheetIndex int) int {
	var a T
	rType := reflect.TypeOf(a)

	colID := 0
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		header := field.Tag.Get("header")
		if header == "" {
			continue
		}

		columnLetter := ColumnIndexToExcelLetter(colID + 1)
		_ = xlsx.SetCellValue(xlsx.GetSheetName(sheetIndex), fmt.Sprintf("%s%d", columnLetter, row), header)
		colID++
	}

	return colID
}

func InsertRow(xlsx *excelize.File, data interface{}, row int, sheetIndex int) {
	rVal := reflect.ValueOf(data)

	colID := 0
	for i := 0; i < rVal.NumField(); i++ {
		field := rVal.Field(i)
		if header := rVal.Type().Field(i).Tag.Get("header"); header == "" {
			continue
		}

		columnLetter := ColumnIndexToExcelLetter(colID + 1)

		if field.Kind() == reflect.Bool {
			str := "No"
			if field.Bool() {
				str = "Yes"
			}
			_ = xlsx.SetCellValue(xlsx.GetSheetName(sheetIndex), fmt.Sprintf("%s%d", columnLetter, row), str)
		} else {
			_ = xlsx.SetCellValue(xlsx.GetSheetName(sheetIndex), fmt.Sprintf("%s%d", columnLetter, row), field.Interface())
		}

		colID++
	}
}

func SetColumnAutoWidth(f *excelize.File, sheetName string, data interface{}) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		panic("data must be a slice")
	}

	if val.Len() == 0 {
		return
	}

	elemType := val.Index(0).Type()
	headers := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		if header := elemType.Field(i).Tag.Get("header"); header != "" {
			headers = append(headers, header)
		}
	}

	// Create a 2D slice of strings to store the data
	columnData := [][]string{}
	for r := 0; r < val.Len(); r++ {
		elem := val.Index(r)
		rowData := []string{}
		for c := 0; c < elem.NumField(); c++ {
			if header := elemType.Field(c).Tag.Get("header"); header != "" {
				rowData = append(rowData, formatValue(elem.Field(c)))
			}
		}
		columnData = append(columnData, rowData)
	}

	// Include headers in the column width calculation
	columnData = append([][]string{headers}, columnData...)

	for col := 0; col < len(headers); col++ {
		maxWidth := 0
		for row := 0; row < len(columnData); row++ {
			value := columnData[row][col]
			width := utf8.RuneCountInString(value)
			if width > maxWidth {
				maxWidth = width
			}
		}
		// Adding a little extra space for padding
		maxWidth += 2
		colName, _ := excelize.ColumnNumberToName(col + 1)
		f.SetColWidth(sheetName, colName, colName, float64(maxWidth))
	}
}

func ColumnIndexToExcelLetter(index int) string {
	result := ""
	for index > 0 {
		index--
		result = string(rune('A'+(index%26))) + result
		index /= 26
	}
	return result
}

func formatValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "Yes"
		}
		return "No"
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", v.Float())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// WriteXLSXUploadMultipleSheet used to write error message in excel with multiple sheets
func WriteXLSXUploadMultipleSheetByPath(filePath string, sheetMetas []ErrorFileSheetMeta) (bytes.Buffer, error) {
	stream, err := os.Open(filePath)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer stream.Close()

	xlsx, err := excelize.OpenReader(stream)
	if err != nil {
		return bytes.Buffer{}, err
	}

	var buf bytes.Buffer
	for _, sheetMeta := range sheetMetas {
		groupedMessages := groupValidationErrors(sheetMeta.ErrorValidator.ValidationErrors())

		// Write grouped messages to the specific sheet
		for row, message := range groupedMessages {
			_ = xlsx.SetCellValue(sheetMeta.SheetName, fmt.Sprintf("%s%d", sheetMeta.LastColumnName, row+sheetMeta.StartIndex+1), message)
		}
	}

	if err := xlsx.Write(&buf); err != nil {
		return buf, err
	}

	return buf, nil
}
