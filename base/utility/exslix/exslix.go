package exslix

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"slices"
	"sort"
	"strings"
)

func Export[T any](sheet Sheet[T]) (*excelize.File, error) {
	xlsx := excelize.NewFile()

	if sheetID, err := xlsx.GetSheetIndex(sheet.SheetName); err != nil && sheetID != -1 {
		return nil, err
	} else if sheetID == -1 {
		if _, err := xlsx.NewSheet(sheet.SheetName); err != nil {
			return nil, err
		}
	}

	xlsx.SetActiveSheet(1)
	_ = xlsx.DeleteSheet("Sheet1")

	sw, err := xlsx.NewStreamWriter(sheet.SheetName)
	if err != nil {
		return nil, err
	}

	if err = sw.SetRow("A1",
		[]interface{}{
			excelize.Cell{Value: sheet.Title},
		}); err != nil {
		return nil, err
	}

	if err = sw.SetRow("A2",
		[]interface{}{
			excelize.Cell{Value: sheet.Company},
		}); err != nil {
		return nil, err
	}

	for i, filter := range sheet.Filters {
		styleID, _ := xlsx.NewStyle(filter.Style)
		if err := sw.SetRow(fmt.Sprintf("A%d", 4+i), []interface{}{
			excelize.Cell{Value: filter.Label, StyleID: styleID},
			excelize.Cell{Value: filter.Value, StyleID: styleID},
		}); err != nil {
			return nil, err
		}
	}

	sort.Slice(sheet.Headers, func(i, j int) bool {
		return sheet.Headers[i].Position < sheet.Headers[j].Position
	})

	var (
		labelHeaders        []any
		mapPositionHeaders  = make(map[int]Header)
		mapAttributeHeaders = make(map[string]Header)
		headerAttributes    []string
	)

	for _, header := range sheet.Headers {
		if header.Position < 1 {
			return nil, fmt.Errorf("position in header must be positive number")
		}

		if _, ok := mapPositionHeaders[header.Position]; ok {
			return nil, fmt.Errorf("position in header must be unique")
		}

		if slices.Contains(headerAttributes, header.Attribute) {
			return nil, fmt.Errorf("attribute in header must be unique")
		}

		mapPositionHeaders[header.Position] = header
		mapAttributeHeaders[header.Attribute] = header
		headerAttributes = append(headerAttributes, header.Attribute)

		styleID, _ := xlsx.NewStyle(header.Style)
		labelHeaders = append(labelHeaders, excelize.Cell{Value: header.Label, StyleID: styleID})
	}

	headerRow := len(sheet.Filters) + 5

	if err = sw.SetRow(
		fmt.Sprintf("A%d", headerRow),
		labelHeaders,
	); err != nil {
		return nil, err
	}

	rfVal := reflect.ValueOf(sheet.Values)

	if rfVal.Type().Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("value must be in slice of struct")
	}

	var (
		dt [][]any
	)

	for i := 0; i < rfVal.Len(); i++ {
		var d []any
		for j := 1; j <= len(labelHeaders); j++ {
			styleID, _ := xlsx.NewStyle(mapPositionHeaders[j].ColumnStyle)
			temp := map[string]excelize.Cell{
				mapPositionHeaders[j].Attribute: {
					StyleID: styleID,
				},
			}

			for k := 0; k < rfVal.Index(i).NumField(); k++ {
				rfValueStruct := rfVal.Index(i).Field(k)
				if !strings.Contains(rfValueStruct.Type().String(), "exslix.Value") {
					return nil, fmt.Errorf("struct must use exslix.Value type")
				}

				if rfValueStruct.FieldByName("Attribute").
					String() !=
					mapPositionHeaders[j].Attribute {
					continue
				}

				styleID, _ := xlsx.NewStyle(rfValueStruct.
					FieldByName("Style").
					Interface().(*excelize.Style))

				cell := temp[mapPositionHeaders[j].Attribute]
				cell.Value = rfValueStruct.
					FieldByName("Value").
					Interface()

				if styleID != 0 {
					cell.StyleID = styleID
				}

				temp[mapPositionHeaders[j].Attribute] = cell
			}
			d = append(d, temp[mapPositionHeaders[j].Attribute])
		}
		dt = append(dt, d)
	}

	for i, d := range dt {
		if err = sw.SetRow(
			fmt.Sprintf("A%d", headerRow+(i+1)),
			d,
		); err != nil {
			return nil, err
		}
	}

	return xlsx, nil
}
