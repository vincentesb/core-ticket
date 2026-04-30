package export_helper

import (
	"bytes"
	"core-ticket/base/helpers/base_helper"
	"core-ticket/base/helpers/date_time_helper"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/export_helper/cell_type_constants"
	"core-ticket/base/helpers/export_helper/text_alignment_constants"
	"core-ticket/base/helpers/s3_client_helper"
	"core-ticket/constants/error_code"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// Excel Properties
const (
	DocPropertiesAuthor             = "ESB Reporting" // Set default author in doc propertiesfor exported data
	ReportSheetName                 = "Report"        // Set Sheet name for exported excel
	TransactionReportSheetName      = "Worksheet"     // Set Sheet name for exported excel (type transaction)
	TemplateSheetName               = "Template"      // Set Sheet name for exported excel (template)
	TotalExportHeaderRow            = 6               // Set total header for exported data
	TransactionTotalExportHeaderRow = 3               // Set total header for exported data
	TotalExportTemplateHeaderRow    = 1               // Set total header for exported data
	ExportFontFamily                = "Arial"         // Set default font family for exported data
	ExportFontSize                  = 10              // Set default font size for exported data
	ExportCellWidth                 = 24              // Set cell width in excel for export
)

// Format
// See the documentation: https://xuri.me/excelize/en/style.html#number_format
const (
	GeneralExcelNumberFormat = "General"
	DecimalExcelNumberFormat = "#,##0.00"
)

// Color
const (
	COLOR_BLACK = "#000000"
)

// Border Style
const (
	BORDER_STYLE_LINE = 1
	BORDER_STYLE_DOT  = 4
)

var (
	DocPropertiesDescription = "%s generated from ESB FNB Application" // Set default description in doc properties for exported data
)

type ExportResponse struct {
	Path string `json:"path"`
}

/*
Fields:
- SheetName: sheet name to be exported
- CompanyName: current company name used for document properties
- Identity: current logged identity
- Title: filled with title of exported data. If IsTransaction is false, this fields will also used for filename
- IsTransaction: filled with true if exported data is transaction type. This fields will affect the naming of file
- TransactionName: filled with transaction name when IsTransaction is true
- RefNum: filled with references number when IsTransaction is true
- Header: filled with exported header data. Plain format
- Data: filled with list of data to be exported
- AdvanceHeader: filled with exported advance header data. Define the row data cell value type in this property
- IsAdvanceHeader: must filled with true if you want use advance header. If it set as false, then advance header will be ignored
*/
type ExportRequest struct {
	SheetName       string
	CompanyName     string
	Identity        base_helper.Identity
	Title           string
	IsTransaction   bool
	IsTemplate      bool
	TransactionName string
	RefNum          string
	Header          []string
	Notes           []string
	Filter          []FilterRequest
	// Data can be filled with slice of struct, slice of interface, and map
	// Note that if using map, the printed data may not be sorted by requested order
	Data                    interface{}
	Summary                 interface{}
	FileName                string
	AdvanceHeader           []Header
	IsAdvanceHeader         bool
	UsePrivateReportStorage bool
}

// Header used to set excel header with additional property
type Header struct {
	Label               string
	Type                cell_type_constants.Type            // Default cell type will be General
	Width               float64                             // Default cell width will be equal to ExportCellWidth
	HorizontalAlignment text_alignment_constants.HAlignment // See constants text_alignment_constants.HAlignment for more detail
}

type FilterRequest struct {
	Key   string
	Value string
}

type CellStyleProperties struct {
	BorderStyleId       int
	CustomNumberFormat  string
	HorizontalAlignment string
}

// SetDocProperties used to set excel document properties
func SetDocProperties(f *excelize.File, companyName string, title string, author string, sheetName string, description string) (*excelize.File, error) {
	var err error
	// Set company name in properties
	if err = f.SetAppProps(&excelize.AppProperties{
		Company: companyName,
	}); err != nil {
		return f, err
	}

	// Set properties in the Excel file
	if err = f.SetDocProps(&excelize.DocProperties{
		Creator:     author,
		Title:       title,
		Description: description,
	}); err != nil {
		return f, err
	}

	if err = f.SetSheetName("Sheet1", sheetName); err != nil {
		return f, err
	}

	return f, nil
}

// InitializeExcelToExport used to initialize new file excel with its properties
// This function sets properties and header for exported file
func InitializeExcelToExport(f *excelize.File, request ExportRequest) (*excelize.File, *excelize.StreamWriter, error) {
	var err error
	now := time.Now()

	noOfColumn := 0
	if request.IsAdvanceHeader {
		noOfColumn = len(request.AdvanceHeader)
	} else {
		noOfColumn = len(request.Header)
	}

	sw, err := f.NewStreamWriter(request.SheetName)
	if err != nil {
		return f, nil, err
	}

	styleA1, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font:      &excelize.Font{Bold: true, Family: ExportFontFamily, Size: 16},
	})
	if err != nil {
		return f, nil, err
	}

	if request.IsAdvanceHeader {
		cellIdx := 1
		for _, header := range request.AdvanceHeader {
			width := float64(ExportCellWidth)
			if header.Width > 0 {
				width = header.Width
			}
			if err = sw.SetColWidth(cellIdx, cellIdx, width); err != nil {
				return f, nil, err
			}

			cellIdx++
		}
	} else {
		if err = sw.SetColWidth(1, noOfColumn, ExportCellWidth); err != nil {
			return f, nil, err
		}
	}

	if err = sw.SetRow("A1",
		[]interface{}{
			excelize.Cell{StyleID: styleA1, Value: request.Title},
		}); err != nil {
		return f, nil, err
	}

	styleA2, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font:      &excelize.Font{Bold: true, Family: ExportFontFamily, Size: 10},
	})
	if err != nil {
		return f, nil, err
	}

	if err = sw.SetRow("A2",
		[]interface{}{
			excelize.Cell{StyleID: styleA2, Value: request.CompanyName},
		}); err != nil {
		return f, nil, err
	}

	if err = sw.SetRow("A4",
		[]interface{}{
			excelize.Cell{StyleID: styleA2, Value: "Generated"},
			excelize.Cell{StyleID: styleA2, Value: date_time_helper.FormatDatetimeToString(now, date_time_helper.FORMAT_DATE_TIME_FIRST)},
		}); err != nil {
		return f, nil, err
	}
	return f, sw, nil
}

// InitializeExcelToExportTransaction used to initialize new file excel with its properties
// This function sets properties and header for exported file
func InitializeExcelToExportTransaction(f *excelize.File, request ExportRequest) (*excelize.File, *excelize.StreamWriter, error) {
	var err error

	sw, err := f.NewStreamWriter(request.SheetName)
	if err != nil {
		return f, nil, err
	}

	styleA2, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font:      &excelize.Font{Bold: true, Family: ExportFontFamily, Size: 10},
	})
	if err != nil {
		return f, nil, err
	}

	headerTitle := "Export " + request.TransactionName + " - " + request.RefNum

	if err = sw.SetRow("A1",
		[]interface{}{
			excelize.Cell{StyleID: styleA2, Value: headerTitle},
		}); err != nil {
		return f, nil, err
	}

	return f, sw, nil
}

// InitializeExcelToExportTemplate used to initialize new file excel with its properties
// This function sets properties and header for exported file
func InitializeExcelToExportTemplate(f *excelize.File, request ExportRequest) (*excelize.File, *excelize.StreamWriter, error) {
	var err error

	noOfColumn := 0
	if request.IsAdvanceHeader {
		noOfColumn = len(request.AdvanceHeader)
	} else {
		noOfColumn = len(request.Header)
	}

	sw, err := f.NewStreamWriter(request.SheetName)
	if err != nil {
		return f, nil, err
	}

	if request.IsAdvanceHeader {
		cellIdx := 1
		for _, header := range request.AdvanceHeader {
			width := float64(ExportCellWidth)
			if header.Width > 0 {
				width = header.Width
			}
			if err = sw.SetColWidth(cellIdx, cellIdx, width); err != nil {
				return f, nil, err
			}

			cellIdx++
		}
	} else {
		if err = sw.SetColWidth(1, noOfColumn, ExportCellWidth); err != nil {
			return f, nil, err
		}
	}

	return f, sw, nil
}

// SetHeaderRowStyle returns styles for Row Header data
func SetHeaderRowStyle(f *excelize.File, property CellStyleProperties) (int, error) {
	stylePattern := excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   ExportFontSize,
			Family: ExportFontFamily,
		},
		Border: []excelize.Border{
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "top",
			},
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "right",
			},
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "bottom",
			},
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "left",
			},
		},
	}

	if property.CustomNumberFormat != "" {
		stylePattern.CustomNumFmt = &property.CustomNumberFormat
	}

	if property.HorizontalAlignment != "" {
		stylePattern.Alignment = &excelize.Alignment{
			Horizontal: property.HorizontalAlignment,
		}
	}

	style, err := f.NewStyle(&stylePattern)
	return style, err
}

// SetBodyRowStyle returns styles for Row Body data
func SetBodyRowStyle(f *excelize.File, property CellStyleProperties) (int, error) {
	stylePattern := excelize.Style{
		Font: &excelize.Font{
			Size:   ExportFontSize,
			Family: ExportFontFamily,
		},
		Border: []excelize.Border{
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "top",
			},
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "right",
			},
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "bottom",
			},
			{
				Style: property.BorderStyleId,
				Color: COLOR_BLACK,
				Type:  "left",
			},
		},
	}

	if property.CustomNumberFormat != "" {
		stylePattern.CustomNumFmt = &property.CustomNumberFormat
	}

	if property.HorizontalAlignment != "" {
		stylePattern.Alignment = &excelize.Alignment{
			Horizontal: property.HorizontalAlignment,
		}
	}

	style, err := f.NewStyle(&stylePattern)
	return style, err
}

// Export used to create generate exported data
func Export(request ExportRequest) (xlsx *excelize.File, err error) {
	styleBorder := BORDER_STYLE_DOT
	totalHeaderRow := TotalExportHeaderRow
	if request.IsTransaction {
		if request.SheetName == "" {
			request.SheetName = TransactionReportSheetName
		}
		styleBorder = BORDER_STYLE_LINE
		totalHeaderRow = TransactionTotalExportHeaderRow
	} else if request.IsTemplate {
		if request.SheetName == "" {
			request.SheetName = TemplateSheetName
		}
		totalHeaderRow = TotalExportTemplateHeaderRow
	} else {
		if request.SheetName == "" {
			request.SheetName = ReportSheetName
		}
	}

	// Create new file
	xlsx = excelize.NewFile()
	var sw *excelize.StreamWriter

	if request.IsTransaction {
		xlsx, err = SetDocProperties(xlsx, request.CompanyName, request.Title, request.Identity.Username, request.SheetName, "Export "+request.Title)
		if err != nil {
			return nil, err
		}
		xlsx, sw, err = InitializeExcelToExportTransaction(xlsx, request)

	} else if request.IsTemplate {
		xlsx, err = SetDocProperties(xlsx, request.CompanyName, request.Title, DocPropertiesAuthor, request.SheetName, fmt.Sprintf(DocPropertiesDescription, request.Title))
		if err != nil {
			return nil, err
		}
		xlsx, sw, err = InitializeExcelToExportTemplate(xlsx, request)
	} else {
		xlsx, err = SetDocProperties(xlsx, request.CompanyName, request.Title, DocPropertiesAuthor, request.SheetName, fmt.Sprintf(DocPropertiesDescription, request.Title))
		if err != nil {
			return nil, err
		}
		xlsx, sw, err = InitializeExcelToExport(xlsx, request)

	}
	if err != nil {
		return nil, err
	}

	// Print filtered row to excel (after header, before row data)
	if len(request.Filter) > 0 {
		styleA2, err := xlsx.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "left"},
			Font:      &excelize.Font{Bold: true, Family: ExportFontFamily, Size: 10},
		})
		if err != nil {
			return nil, err
		}
		// Rmove empty space after header
		totalHeaderRow -= 1
		// Set filter data if exists
		for _, item := range request.Filter {
			setColHeader := "A" + strconv.Itoa(totalHeaderRow)
			filterRow := make([]interface{}, 0)
			filterRow = append(filterRow, excelize.Cell{StyleID: styleA2, Value: item.Key})
			filterRow = append(filterRow, excelize.Cell{StyleID: styleA2, Value: item.Value})
			if err = sw.SetRow(setColHeader, filterRow); err != nil {
				return nil, err
			}
			totalHeaderRow += 1
		}
		// Add empty space after filter row
		totalHeaderRow += 1
	}

	// Initialize header style
	styleA6, err := SetHeaderRowStyle(xlsx, CellStyleProperties{
		BorderStyleId:       styleBorder,
		CustomNumberFormat:  GeneralExcelNumberFormat,
		HorizontalAlignment: string(text_alignment_constants.Center),
	})
	if err != nil {
		return nil, err
	}

	if request.Notes != nil {
		notesData := make([]interface{}, 0)
		// Set notes row data
		for _, notes := range request.Notes {
			notesData = append(notesData, excelize.Cell{StyleID: styleA6, Value: notes})
		}
		// Write notes to row
		setColNotes := "A" + strconv.Itoa(totalHeaderRow)
		if err = sw.SetRow(setColNotes, notesData); err != nil {
			return nil, err
		}
		totalHeaderRow++
	}

	headerData := make([]interface{}, 0)

	// Set header row data
	if request.IsAdvanceHeader {
		for _, header := range request.AdvanceHeader {
			headerData = append(headerData, excelize.Cell{StyleID: styleA6, Value: header.Label})
		}
	} else {
		for _, header := range request.Header {
			headerData = append(headerData, excelize.Cell{StyleID: styleA6, Value: header})
		}
	}

	// Write header to row
	setColHeader := "A" + strconv.Itoa(totalHeaderRow)
	if err = sw.SetRow(setColHeader, headerData); err != nil {
		return nil, err
	}

	// Set row data
	lastRowIdx := 0
	if request.IsAdvanceHeader {
		sw, lastRowIdx, err = SetAdvanceRowData(xlsx, sw, request, totalHeaderRow, styleBorder)
		if err != nil {
			return nil, err
		}
	} else {
		sw, lastRowIdx, err = SetRowData(xlsx, sw, request, totalHeaderRow, styleBorder)
		if err != nil {
			return nil, err
		}
	}

	if request.Summary != nil {
		// Set row summary
		sw, err = SetSummaryRowData(xlsx, sw, request.Summary, lastRowIdx, styleBorder)
		if err != nil {
			return nil, err
		}
	}

	if err = sw.Flush(); err != nil {
		return nil, err
	}

	return xlsx, nil
}

// ExportData used to create generate exported data and upload to S3
func ExportData(S3Client s3_client_helper.S3Client, request ExportRequest) (response ExportResponse, err error) {
	xlsx, err := Export(request)
	defer func() {
		if xlsx != nil {
			_ = xlsx.Close()
		}
	}()

	if err != nil {
		return response, err
	}

	response, err = UploadToS3(S3Client, xlsx, request)
	if err != nil {
		return response, err
	}

	return response, nil
}

func GenerateFileName(request ExportRequest) string {
	var fileName string
	randUUID, _ := uuid.NewRandom()
	if request.FileName != "" {
		fileName = fmt.Sprintf("%s.xlsx", request.FileName)
	} else if request.IsTransaction {
		fileName = fmt.Sprintf("Export-%s-%s-%s.xlsx", request.TransactionName, request.RefNum, randUUID.String())
	} else if request.IsTemplate {
		fileName = fmt.Sprintf("Template-%s-%s-%s.xlsx", request.Identity.Username, request.Title, randUUID.String())
	} else {
		fileName = fmt.Sprintf("Export-%s-%s-%s.xlsx", request.Identity.Username, request.Title, randUUID.String())
	}

	return fileName
}

func UploadToS3(S3Client s3_client_helper.S3Client, xlsx *excelize.File, request ExportRequest) (ExportResponse, error) {
	var buf bytes.Buffer
	if err := xlsx.Write(&buf); err != nil {
		return ExportResponse{}, err
	}

	// Set filename
	fileName := GenerateFileName(request)
	key := fmt.Sprintf("export/%s/%s", request.Identity.CompanyCode, fileName)

	reader := bytes.NewReader(buf.Bytes()) // implements io.ReadSeeker

	var (
		path string
		err  error
	)

	if request.UsePrivateReportStorage {
		path, err = S3Client.UploadReport(reader, key)
	} else {
		path, err = S3Client.UploadRaw(reader, key)
	}

	if err != nil {
		return ExportResponse{}, error_helper.New(errors.New("failed to upload exported data"), error_code.UnknownError)
	}

	return ExportResponse{Path: path}, nil
}

func SetRowData(xlsx *excelize.File, sw *excelize.StreamWriter, request ExportRequest, totalHeaderRow int, style int) (*excelize.StreamWriter, int, error) {
	dataValue := reflect.ValueOf(request.Data)
	lastRowIdx := totalHeaderRow

	if dataValue.Kind() != reflect.Slice {
		return sw, lastRowIdx, error_helper.New(errors.New("failed to extract data"), error_code.UnknownError)
	}

	length := dataValue.Len()
	if length == 0 {
		return sw, lastRowIdx, nil
	}

	generalFormatStyle, decimalFormatStyle, err := defineStyles(xlsx, style)
	if err != nil {
		return sw, lastRowIdx, nil
	}

	getStyleByType := getStyleFunc(generalFormatStyle, decimalFormatStyle)

	// Loop for each row in data
	for i := 0; i < length; i++ {
		data := dataValue.Index(i).Interface()
		row := make([]interface{}, 0)
		t := reflect.TypeOf(data)
		elem := dataValue.Index(i)

		if elem.Kind() == reflect.Map {
			for _, key := range elem.MapKeys() {
				fieldValue := elem.MapIndex(key).Interface()
				styleID := getStyleByType(fieldValue)
				row = append(row, excelize.Cell{StyleID: styleID, Value: fieldValue})
			}
		} else if elem.Kind() == reflect.Slice {
			// Handle []interface{}
			for j := 0; j < elem.Len(); j++ {
				sliceElem := elem.Index(j)
				rawValue := sliceElem.Interface()
				fieldValue := normalizeExcelValue(rawValue)
				styleID := getStyleByType(fieldValue)
				row = append(row, excelize.Cell{StyleID: styleID, Value: fieldValue})
			}
		} else {
			numFields := t.NumField()
			// Loop for each col data in a row
			for j := 0; j < numFields; j++ {
				rawValue := reflect.ValueOf(data).Field(j).Interface()
				fieldValue := normalizeExcelValue(rawValue)
				styleID := getStyleByType(fieldValue)
				row = append(row, excelize.Cell{StyleID: styleID, Value: fieldValue})
			}
		}

		// Set row position
		colIdx := i + totalHeaderRow + 1
		cell, errCell := excelize.CoordinatesToCellName(1, colIdx)
		if errCell != nil {
			return sw, lastRowIdx, errCell
		}

		// Write to excel row
		if errRow := sw.SetRow(cell, row); errRow != nil {
			return sw, lastRowIdx, errRow
		}

		lastRowIdx = colIdx
	}

	return sw, lastRowIdx, nil
}

func SetAdvanceRowData(xlsx *excelize.File, sw *excelize.StreamWriter, request ExportRequest, totalHeaderRow int, borderStyle int) (*excelize.StreamWriter, int, error) {
	dataValue := reflect.ValueOf(request.Data)
	lastRowIdx := totalHeaderRow

	if dataValue.Kind() != reflect.Slice {
		return sw, lastRowIdx, error_helper.New(errors.New("failed to extract data"), error_code.UnknownError)
	}

	// Initialize column types and styles
	columnTypes, err := initializeColumnStyleFormat(xlsx, extractHeaderTypes(request.AdvanceHeader), borderStyle)
	if err != nil {
		return sw, 0, err
	}

	length := dataValue.Len()
	if length == 0 {
		return sw, lastRowIdx, nil
	}

	// Loop for each row in data
	for i := 0; i < length; i++ {
		data := dataValue.Index(i).Interface()
		row := make([]interface{}, 0)
		t := reflect.TypeOf(data)
		numFields := t.NumField()

		// Loop for each col data in a row
		for j := 0; j < numFields; j++ {
			styleId := columnTypes[j].DefaultCellStyle
			fieldValue := reflect.ValueOf(data).Field(j).Interface()
			row = append(row, excelize.Cell{StyleID: styleId, Value: fieldValue})
		}

		// Set row position
		colIdx := i + totalHeaderRow + 1
		cell, errCell := excelize.CoordinatesToCellName(1, colIdx)
		if errCell != nil {
			return sw, lastRowIdx, errCell
		}

		// Write to excel row
		if errRow := sw.SetRow(cell, row); errRow != nil {
			return sw, lastRowIdx, errRow
		}

		lastRowIdx = colIdx
	}

	return sw, lastRowIdx, nil
}

func SetSummaryRowData(xlsx *excelize.File, sw *excelize.StreamWriter, data interface{}, lastRowIdx int, style int) (*excelize.StreamWriter, error) {
	dataValue := reflect.ValueOf(data)

	if dataValue.Kind() != reflect.Slice {
		return sw, error_helper.New(errors.New("failed to extract data"), error_code.UnknownError)
	}

	length := dataValue.Len()
	if length == 0 {
		return sw, nil
	}

	generalFormatStyle, decimalFormatStyle, err := defineStyles(xlsx, style)
	if err != nil {
		return sw, nil
	}

	getStyleByType := getStyleFunc(generalFormatStyle, decimalFormatStyle)

	// Loop for each row in data
	for i := 0; i < length; i++ {
		data := dataValue.Index(i).Interface()
		row := make([]interface{}, 0)
		t := reflect.TypeOf(data)
		numFields := t.NumField()

		// Loop for each col data in a row
		for j := 0; j < numFields; j++ {
			rawValue := reflect.ValueOf(data).Field(j).Interface()
			fieldValue := normalizeExcelValue(rawValue)
			styleID := getStyleByType(fieldValue)
			row = append(row, excelize.Cell{StyleID: styleID, Value: fieldValue})
		}

		// Set row position
		colIdx := i + lastRowIdx + 1
		cell, errCell := excelize.CoordinatesToCellName(1, colIdx)
		if errCell != nil {
			return sw, errCell
		}

		// Write to excel row
		if errRow := sw.SetRow(cell, row); errRow != nil {
			return sw, errRow
		}
	}

	return sw, nil
}

// Number format section
func defineStyles(xlsx *excelize.File, style int) (int, int, error) {
	generalFormatStyle, err := SetBodyRowStyle(xlsx, CellStyleProperties{
		BorderStyleId:      style,
		CustomNumberFormat: GeneralExcelNumberFormat,
	})
	if err != nil {
		return 0, 0, err
	}

	decimalFormatStyle, err := SetBodyRowStyle(xlsx, CellStyleProperties{
		BorderStyleId:      style,
		CustomNumberFormat: DecimalExcelNumberFormat,
	})
	if err != nil {
		return 0, 0, err
	}

	return generalFormatStyle, decimalFormatStyle, nil
}

func getStyleFunc(generalFormatStyle, decimalFormatStyle int) func(interface{}) int {
	styles := map[reflect.Kind]int{
		reflect.Float64: decimalFormatStyle,
		reflect.Int:     generalFormatStyle,
		reflect.String:  generalFormatStyle,
		// Add more types as needed
	}

	return func(value interface{}) int {
		if styleID, ok := styles[reflect.TypeOf(value).Kind()]; ok {
			return styleID
		}
		return generalFormatStyle // default to generalFormatStyle
	}
}

// normalizeExcelValue converts various Go types into values safe for Excel cells.
//
// This helper prevents issues like printing nil pointers as memory addresses
// and ensures consistent formatting for exported data.
//
// Rules:
//   - nil → ""
//   - *float64 → dereferenced value or "-"
//   - Numeric types → returned as-is
//   - string → returned as-is
//   - Other types → converted to string using fmt.Sprintf("%v")
func normalizeExcelValue(v interface{}) interface{} {
	switch val := v.(type) {
	case nil:
		return ""
	case *float64:
		if val != nil {
			return *val
		}
		return "-"
	case float64, float32, int, int64, int32:
		return val
	case string:
		return val
	case *string:
		if val != nil {
			return *val
		}
		return "-"
	default:
		return fmt.Sprintf("%v", val)
	}
}
