package s3_client_helper

import (
	"bytes"
	"core-ticket/base/helpers/base_helper"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type S3Client interface {
	BulkDelete(fileUrls []string) error
	Delete(fileUrl string) (bool, error)
	DownloadTemp(fileUrl string) (string, error)
	Upload(fileHeader *multipart.FileHeader, destinationFilePath string) (string, error)
	UploadRaw(file io.ReadSeeker, destinationFilePath string) (string, error)
	UploadReport(file io.ReadSeeker, destinationFilePath string) (string, error)
}

type S3ClientImpl struct {
	// Default S3 (from S3_* env vars)
	Client   *s3.S3
	Config   *S3Config
	Uploader *s3manager.Uploader

	// Optional Report S3 (from REPORT_S3_* env vars). If nil, fallback to Default fields above.
	ReportClient   *s3.S3
	ReportConfig   *S3Config
	ReportUploader *s3manager.Uploader
}

/*
NewS3Client creates a new S3Client instance with the configuration values obtained from environment variables using the getConfig function. It initializes a new AWS session with the specified region, endpoint, access key, and secret key. If any error occurs during session creation, it panics with the error message.

Returns:
- S3Client: A new S3Client instance with the initialized AWS S3 client, configuration, and uploader.

Panic:
- If an error occurs during AWS session creation, it panics with the error message.
*/
func NewS3Client() S3Client {
	config := getConfig()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Endpoint:    aws.String(config.EndPoint),
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}

	svc := &S3ClientImpl{
		Client:   s3.New(sess),
		Config:   config,
		Uploader: s3manager.NewUploader(sess),
	}

	// Initialize report uploader/config if REPORT_S3_ENDPOINT is present
	if rc := getReportConfig(); rc != nil {
		repSess, err := session.NewSession(&aws.Config{
			Region:      aws.String(rc.Region),
			Endpoint:    aws.String(rc.EndPoint),
			Credentials: credentials.NewStaticCredentials(rc.AccessKey, rc.SecretKey, ""),
		})
		if err == nil {
			svc.ReportConfig = rc
			svc.ReportClient = s3.New(repSess)
			svc.ReportUploader = s3manager.NewUploader(repSess)
		}
	}

	return svc
}

/*
BulkDelete deletes multiple objects from the S3 bucket based on the provided list of file URLs.

Parameters:
- fileUrls: a slice of strings representing the URLs of the files to be deleted from the S3 bucket.

Returns:
- error: an error if any occurred during the deletion process, nil otherwise.
*/
func (svc *S3ClientImpl) BulkDelete(fileUrls []string) error {
	var keys []*s3.ObjectIdentifier
	for _, fileUrl := range fileUrls {
		key, err := svc.getKeyFromUrl(fileUrl)
		if err != nil {
			return err
		}
		keys = append(keys, &s3.ObjectIdentifier{Key: aws.String(key)})
	}

	if len(keys) == 0 {
		return nil
	}

	_, err := svc.Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(svc.Config.Bucket),
		Delete: &s3.Delete{Objects: keys},
	})
	return err
}

/*
Delete deletes the object from the S3 bucket based on the provided file URL.

Parameters:
- fileUrl (string): The URL of the file to be deleted from the S3 bucket.

Returns:
- bool: true if the object is successfully deleted, false otherwise.
- error: An error if any occurred during the deletion process.
*/
func (svc *S3ClientImpl) Delete(fileUrl string) (bool, error) {
	key, err := svc.getKeyFromUrl(fileUrl)
	if err != nil {
		return false, err
	}

	_, err = svc.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(svc.Config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

/*
Upload uploads a file to the specified destinationFilePath in the S3 bucket configured in the S3ClientImpl instance.

Parameters:
- fileHeader: The multipart file header of the file to be uploaded.
- destinationFilePath: The path where the file will be uploaded in the S3 bucket.

Returns:
- string: The URL of the uploaded file.
- error: An error if the upload operation fails.
*/
func (svc *S3ClientImpl) Upload(fileHeader *multipart.FileHeader, destinationFilePath string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	return svc.uploadToS3(file, destinationFilePath)
}

/*
UploadRaw uploads a file to the specified destinationFilePath in the S3 bucket configured in the S3ClientImpl instance.
It sets the content type based on the file extension and uploads the file with public-read ACL.
Returns the URL of the uploaded file on success, or an error if the upload fails.

Parameters:
- file: io.ReadSeeker - The file to upload.
- destinationFilePath: string - The path where the file will be stored in the S3 bucket.

Returns:
- string: The URL of the uploaded file.
- error: An error if the upload operation fails.
*/
func (svc *S3ClientImpl) UploadRaw(file io.ReadSeeker, destinationFilePath string) (string, error) {
	return svc.uploadToS3(file, destinationFilePath)
}

// uploadToS3 is a helper function to upload files to S3 with public-read ACL.
func (svc *S3ClientImpl) uploadToS3(file io.ReadSeeker, destinationFilePath string) (string, error) {
	contentType := mime.TypeByExtension(filepath.Ext(destinationFilePath))

	uploadResult, err := svc.Uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Body:        file,
		Bucket:      aws.String(svc.Config.Bucket),
		Key:         aws.String(destinationFilePath),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// Ensure ACL is explicitly set to public-read
	_, err = svc.Uploader.S3.PutObjectAcl(&s3.PutObjectAclInput{
		Bucket: aws.String(svc.Config.Bucket),
		Key:    aws.String(destinationFilePath),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	return uploadResult.Location, nil
}

/*
getKeyFromUrl extracts the key from the given S3 URL.

Parameters:
- s3url (string): The S3 URL from which to extract the key.

Returns:
- string: The extracted key.
- error: An error if the URL parsing fails.

Example:

	key, err := svc.getKeyFromUrl("s3://bucket-name/path/to/file.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(key) // Output: path/to/file.txt
*/
func (svc *S3ClientImpl) getKeyFromUrl(s3url string) (string, error) {
	parsedURL, err := url.Parse(s3url)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(parsedURL.Path, "/"), nil
}

/*
DownloadTemp downloads the file from S3 and saves it to a temporary file.
You should delete (manually) the file after you're done using it.

Parameters:
- fileUrl: string - The URL of the file in the S3 bucket.

Returns:
- string: The path to the temporary file where the downloaded content is stored.
- error: An error if the download operation fails.
*/
func (svc *S3ClientImpl) DownloadTemp(fileUrl string) (string, error) {
	fileKey, err := svc.getKeyFromUrl(fileUrl)
	if err != nil {
		return "", fmt.Errorf("error extracting file key: %v", err)
	}

	downloadedFile, err := svc.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(svc.Config.Bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", fmt.Errorf("error downloading file from S3: %v", err)
	}
	defer downloadedFile.Body.Close()

	// Read the file content into memory
	fileBuffer := &bytes.Buffer{}
	_, err = io.Copy(fileBuffer, downloadedFile.Body)
	if err != nil {
		return "", fmt.Errorf("error reading file content: %v", err)
	}

	// Create a temporary file to store the downloaded content
	tempFile, err := os.CreateTemp("", fmt.Sprintf("downloaded-%d-*", time.Now().UnixNano()))
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}

	// Write the content into the temporary file
	_, err = tempFile.Write(fileBuffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("error writing to temp file: %v", err)
	}

	return tempFile.Name(), nil
}

func GenerateDestinationPath(identity base_helper.Identity, title, fileName string) string {
	randomUUID, _ := uuid.NewV7()
	ext := filepath.Ext(fileName)
	finalFileName := fmt.Sprintf("%s-%s-%s%s", title, identity.Username, randomUUID.String(), ext)
	filePath := fmt.Sprintf("export/%s/%s", identity.CompanyCode, finalFileName)

	return filePath
}

// UploadReport uploads a file to the report S3 configuration
func (svc *S3ClientImpl) UploadReport(file io.ReadSeeker, destinationFilePath string) (string, error) {
	if svc.ReportUploader != nil && svc.ReportConfig != nil {
		contentType := mime.TypeByExtension(filepath.Ext(destinationFilePath))
		uploadResult, err := svc.ReportUploader.Upload(&s3manager.UploadInput{
			Body:        file,
			Bucket:      aws.String(svc.ReportConfig.Bucket),
			Key:         aws.String(destinationFilePath),
			ContentType: aws.String(contentType),
		})
		if err != nil {
			return "", err
		}

		return uploadResult.Location, nil
	}

	// Fallback to default uploader
	return svc.uploadToS3(file, destinationFilePath)
}
