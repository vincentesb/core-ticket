package attachment_helper

import (
	"core-ticket/base/helpers/error_helper"
	"core-ticket/constants/error_code"
	"core-ticket/constants/error_message"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
)

// Default value
const (
	megaBytes           = 1024 * 1024
	TotalMaxFileSize    = 10 * megaBytes // 10 mb
	ErrMsgTotalFileSize = "10mb"         // for size 10 mb
	MaxFileSize         = 2 * megaBytes  // 2 mb
	ErrMsgFileSize      = "2mb"          // for size 2 mb
)

// SupportedFileExt supported file extensions for attachment
var SupportedFileExt = []string{".pdf", ".bmp", ".gif", ".jpe", ".jpeg", ".jpg", ".png"}

/*
GetAttachmentFiles extracts files and existingFileUrls from a multipart-form request.

Parameters:
- form: Pointer to a multipart form.

Returns:
- files: Array of file headers (SupportedFileExt).
- existingFiles: Array of existing file URLs.
- err: Error encountered (file size or file type validation).
*/
func GetAttachmentFiles(form *multipart.Form) (files []*multipart.FileHeader, existingFiles []string, err error) {
	var totalSize int64
	for key, value := range form.File {
		if !strings.Contains(key, "files") {
			continue
		}
		file := value[0]
		if file.Size > MaxFileSize {
			return nil, nil, error_helper.New(errors.New(error_message.ErrFileTooBig(ErrMsgFileSize)), error_code.ValidationError)
		}
		fileExt := GetFileExtension(file)
		if !slices.Contains(SupportedFileExt, fileExt) {
			return nil, nil, error_helper.New(errors.New(error_message.ErrFileTypeNotValid), error_code.ValidationError)
		}
		totalSize += file.Size
		files = append(files, file)
	}
	// files should be below TotalMaxFileSize
	if totalSize > TotalMaxFileSize {
		return nil, nil, error_helper.New(errors.New(fmt.Sprintf("Total files too big, maximum file size is %s", ErrMsgTotalFileSize)), error_code.ValidationError)
	}
	for key, value := range form.Value {
		if !strings.Contains(key, "existing") {
			continue
		}
		existingFiles = append(existingFiles, value[0])
	}
	return files, existingFiles, nil
}

/*
GetFileExtension function to retrieve file path extension
*/
func GetFileExtension(file *multipart.FileHeader) string {
	return strings.ToLower(filepath.Ext(file.Filename))
}
