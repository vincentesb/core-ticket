package excel_reader

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
)

type ExcelFile struct {
	File *multipart.FileHeader `form:"file"`
}

func ExtractExcelFile[T interface{}](file *multipart.FileHeader, opts ...excelize.Options) ([]T, error) {
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

	var rowDataList []T
	for _, sheet := range xlsx.GetSheetList() {
		rows, err := xlsx.GetRows(sheet)
		if err != nil {
			return nil, fmt.Errorf("error getting rows from sheet %s: %v", sheet, err)
		}
		for i, row := range rows {
			if i > 0 {
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
	}
	return rowDataList, nil
}
