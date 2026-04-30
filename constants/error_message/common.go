package error_message

import "fmt"

const (
	ErrFailedGetServerCode = "Failed to get server code"
	ErrFailedGetDBName     = "Failed to retrieve Database"
	ErrFailedGetCompanyID  = "Failed to retrieve Company ID"
	ErrUnkownError         = "Unknown error"
	ErrFailedGetUsername   = "Failed to retrieve username"
	ErrFailedGetFile       = "Failed to retrieve file"
	ErrFileTypeNotValid    = "File type not valid"
)

func ErrFileTooBig(maxSize string) string {
	return fmt.Sprintf("File too big, maximum file size is %s", maxSize)
}
