package import_helper

import (
	"bytes"
	"core-ticket/base/helpers/base_helper"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/file_helper/s3_helper/excel_helper"
	"core-ticket/base/helpers/s3_client_helper"
	"core-ticket/constants/error_message"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

// File Content Type
const (
	FILE_CONTENT_TYPE_SHEET = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	FILE_CONTENT_TYPE_EXCEL = "application/vnd.ms-excel"
)

// Default value
const (
	MAX_FILE_SIZE     = 40 * 1024 * 1024 // 40 mb
	ERR_MSG_FILE_SIZE = "40mb"           // for size 40 mb

	MAX_FILE_SIZE_2000KB     = 2000 * 1024 // 2000 KB
	ERR_MSG_FILE_SIZE_2000KB = "2000 KB"   // for size 2000 KB
)

// Upload file default value
const (
	UPLOAD_FILE_SIZE     = 2 * 1024 * 1024 // 2 mb
	UPLOAD_FILE_SIZE_MSG = "2mb"           // for size 2 mb
)

type ImportResponse struct {
	ErrorExcelPathUrl string `json:"url,omitempty"`
}

/*
Fields:
- File: file excel to be uploaded
- S3client: S3Client
- Identity: current logged identity
- ErrorValidator: error from validator
- FeatureName: feature name of uploaded data
- LastColumnName: filled with last column of excel (used to set column of error message). Example: "L" -> column L in excel
- StartIndex: start index row of excel (example: excel header is 1, then start index will be 0)
*/
type UploadImportErrorFileRequest struct {
	File           *multipart.FileHeader
	S3Client       s3_client_helper.S3Client
	Identity       base_helper.Identity
	ErrorValidator *error_helper.Error
	FeatureName    string
	LastColumnName string
	StartIndex     int
}

/*
Fields:
- FilePath: excel temp filepath to be modified and uploaded
- S3client: S3Client
- Identity: current logged identity
- ErrorValidator: error from validator
- FeatureName: feature name of uploaded data
- LastColumnName: filled with last column of excel (used to set column of error message). Example: "L" -> column L in excel
- StartIndex: start index row of excel (example: excel header is 1, then start index will be 0)
*/
type UploadImportErrorFilePathRequest struct {
	FilePath       string
	S3Client       s3_client_helper.S3Client
	Identity       base_helper.Identity
	ErrorValidator *error_helper.Error
	FeatureName    string
	LastColumnName string
	StartIndex     int
}

/*
Fields:
- File: file excel to be uploaded
- S3client: S3Client
- Identity: current logged identity
- FeatureName: feature name of uploaded data
- Sheets: list of Sheets meta
*/
type UploadImportErrorFileRequestMultipleSheet struct {
	File        *multipart.FileHeader
	S3Client    s3_client_helper.S3Client
	Identity    base_helper.Identity
	FeatureName string
	Sheets      []excel_helper.ErrorFileSheetMeta
}

/*
Fields:
- File: file excel to be uploaded
- S3client: S3Client
- Identity: current logged identity
- FeatureName: feature name of uploaded data
- Sheets: list of Sheets meta
*/
type UploadImportErrorFileRequestMultipleSheetByPath struct {
	File        string
	S3Client    s3_client_helper.S3Client
	Identity    base_helper.Identity
	FeatureName string
	Sheets      []excel_helper.ErrorFileSheetMeta
}

/*
FileValidation used to validate uploaded file when importing data (excel file)
Fields:
- maxMemory: maximum size of file (example: 40 * 1024 * 1024) -> 40mb
- errFileMsg: error file message if exceed maxMemory (example: "40mb")
*/
func FileValidation(c *gin.Context, maxMemory int64, errFileMsg string) (file *multipart.FileHeader, errMsg error) {
	err := c.Request.ParseMultipartForm(maxMemory)
	if err != nil {
		errMsg = err
		return
	}

	file, err = c.FormFile("file")
	if err != nil {
		errMsg = errors.New(error_message.ErrFailedGetFile)
		return
	}

	// size limit based on errFileMsg
	if file.Size > maxMemory {
		errMsg = errors.New(error_message.ErrFileTooBig(errFileMsg))
		return
	}

	fileContentType := file.Header.Get("Content-Type")
	// allowed file extension xlsx and xls
	if fileContentType != FILE_CONTENT_TYPE_SHEET && fileContentType != FILE_CONTENT_TYPE_EXCEL {
		errMsg = errors.New(error_message.ErrFileTypeNotValid)
		return
	}

	return
}

// UploadImportErrorFile used to upload excel file that contains error message to S3
func UploadImportErrorFile(request UploadImportErrorFileRequest) (response ImportResponse) {
	buf, _ := excel_helper.WriteXLSXUpload(request.File, request.StartIndex, request.LastColumnName, request.ErrorValidator.ValidationErrors())
	path, _ := request.S3Client.UploadRaw(
		bytes.NewReader(buf.Bytes()),
		fmt.Sprintf("export/%s/%s/%s - Error - Upload %s.xlsx", request.Identity.CompanyCode, request.FeatureName, request.Identity.Username, request.FeatureName),
	)

	response.ErrorExcelPathUrl = path
	return
}

// UploadImportErrorFileByPath used to upload excel file that contains error message to S3, read from temp file path
func UploadImportErrorFileByPath(request UploadImportErrorFilePathRequest) (response ImportResponse) {
	buf, _ := excel_helper.WriteXLSXUploadByPath(request.FilePath, request.StartIndex, request.LastColumnName, request.ErrorValidator.ValidationErrors())
	path, _ := request.S3Client.UploadRaw(
		bytes.NewReader(buf.Bytes()),
		fmt.Sprintf("export/%s/%s/%s - Error - Upload %s.xlsx", request.Identity.CompanyCode, request.FeatureName, request.Identity.Username, request.FeatureName),
	)

	response.ErrorExcelPathUrl = path
	return
}

// UploadImportErrorFile used to upload excel file that contains error message to S3
func UploadImportErrorFileMultipleSheet(request UploadImportErrorFileRequestMultipleSheet) (response ImportResponse) {
	buf, _ := excel_helper.WriteXLSXUploadMultipleSheet(request.File, request.Sheets)
	path, _ := request.S3Client.UploadRaw(
		bytes.NewReader(buf.Bytes()),
		fmt.Sprintf("export/%s/%s/%s - Error - Upload %s.xlsx", request.Identity.CompanyCode, request.FeatureName, request.Identity.Username, request.FeatureName),
	)

	response.ErrorExcelPathUrl = path
	return
}

// UploadImportErrorFile used to upload excel file that contains error message to S3
func UploadImportErrorFileMultipleSheetByPath(request UploadImportErrorFileRequestMultipleSheetByPath) (response ImportResponse) {
	buf, _ := excel_helper.WriteXLSXUploadMultipleSheetByPath(request.File, request.Sheets)
	path, _ := request.S3Client.UploadRaw(
		bytes.NewReader(buf.Bytes()),
		fmt.Sprintf("export/%s/%s/%s - Error - Upload %s.xlsx", request.Identity.CompanyCode, request.FeatureName, request.Identity.Username, request.FeatureName),
	)

	response.ErrorExcelPathUrl = path
	return
}
