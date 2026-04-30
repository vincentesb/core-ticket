package export_helper

import (
	"core-ticket/base/helpers/export_helper/cell_type_constants"
	"core-ticket/base/helpers/export_helper/text_alignment_constants"
	"strings"

	"github.com/xuri/excelize/v2"
)

type HeaderType struct {
	CellType            cell_type_constants.Type
	HorizontalAlignment text_alignment_constants.HAlignment
}

type ColumnType struct {
	NumberFormat     string
	DefaultCellStyle int
}

func getCellNumberFormat(cellFormat cell_type_constants.Type) string {
	switch cellFormat {
	case cell_type_constants.Money:
		cellFormat = "dollar"
	case cell_type_constants.Number:
		cellFormat = "integer"
	case cell_type_constants.String:
		cellFormat = "@"
	case cell_type_constants.Integer:
		cellFormat = "0"
	case cell_type_constants.Date:
		cellFormat = "YYYY-MM-DD"
	case cell_type_constants.DateTime:
		cellFormat = "YYYY-MM-DD HH:MM:SS"
	case cell_type_constants.Price:
		cellFormat = "#,##0.00"
	case cell_type_constants.Dollar:
		cellFormat = "[$$-1009]#,##0.00;[RED]-[$$-1009]#,##0.00"
	case cell_type_constants.Euro:
		cellFormat = "#,##0.00 [$€-407];[RED]-#,##0.00 [$€-407]"
	}

	var escaped strings.Builder
	ignoreUntil := ""
	for i := 0; i < len(cellFormat); i++ {
		c := string(cellFormat[i])
		if ignoreUntil == "" && c == "[" {
			ignoreUntil = "]"
		} else if ignoreUntil == "" && c == "\"" {
			ignoreUntil = "\""
		} else if ignoreUntil == c {
			ignoreUntil = ""
		}
		if ignoreUntil == "" && (c == " " || c == "-" || c == "(" || c == ")") && (i == 0 || string(cellFormat[i-1]) != "_") {
			escaped.WriteString("\\")
			escaped.WriteString(c)
		} else {
			escaped.WriteString(c)
		}
	}

	return escaped.String()
}

func initializeColumnStyleFormat(f *excelize.File, headerTypes []HeaderType, borderStyle int) ([]ColumnType, error) {
	columnTypes := make([]ColumnType, len(headerTypes))
	for i, v := range headerTypes {
		numberFormat := getCellNumberFormat(v.CellType)
		cellStyleIdx, err := SetBodyRowStyle(
			f,
			CellStyleProperties{
				BorderStyleId:       borderStyle,
				CustomNumberFormat:  numberFormat,
				HorizontalAlignment: string(v.HorizontalAlignment),
			},
		)
		if err != nil {
			return nil, err
		}
		columnTypes[i] = ColumnType{
			NumberFormat:     numberFormat,
			DefaultCellStyle: cellStyleIdx,
		}
	}
	return columnTypes, nil
}

func extractHeaderTypes(headers []Header) []HeaderType {
	var types []HeaderType
	for _, header := range headers {
		types = append(types, HeaderType{
			CellType:            header.Type,
			HorizontalAlignment: header.HorizontalAlignment,
		})
	}
	return types
}
