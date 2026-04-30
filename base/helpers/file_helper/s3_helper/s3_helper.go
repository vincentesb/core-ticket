package s3_helper

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

/*
Deprecated: This method is deprecated and will be removed soon.

Note: Please use s3_client_helper.S3Client.UploadRaw instead.
*/
func UploadWithPublicACLRaw(file io.ReadSeeker, path string, ext string) (string, error) {
	contentType := mime.TypeByExtension(ext)

	s3Session, _ := session.NewSession(&aws.Config{
		Region:   aws.String(viper.Get("IMAGE_UPLOAD_REGION").(string)),
		Endpoint: aws.String(fmt.Sprintf("https://%s", viper.Get("IMAGE_UPLOAD_ENDPOINT").(string))),
		Credentials: credentials.NewStaticCredentials(
			viper.Get("IMAGE_UPLOAD_KEY").(string),
			viper.Get("IMAGE_UPLOAD_SECRET").(string),
			"",
		),
	})

	svc := s3.New(s3Session)

	bucket := viper.Get("IMAGE_UPLOAD_BUCKET").(string)
	_, err := svc.PutObject(&s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Body:        file,
		Bucket:      aws.String(bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(path),
	})
	if err != nil {
		return "", err
	}

	endpoint := viper.Get("IMAGE_UPLOAD_ENDPOINT").(string)
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, path), nil
}

/*
Deprecated: This method is deprecated and will be removed soon.

Note: Please use s3_client_helper.S3Client.Upload instead.
*/
func UploadWithPublicACL(fileHeader *multipart.FileHeader, path string, ext string, s3Session *session.Session) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	contentType := mime.TypeByExtension(ext)

	svc := s3.New(s3Session)

	bucket := viper.Get("IMAGE_UPLOAD_BUCKET").(string)
	_, err = svc.PutObject(&s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Body:        file,
		Bucket:      aws.String(bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(path),
	})
	if err != nil {
		return "", err
	}

	endpoint := viper.Get("IMAGE_UPLOAD_ENDPOINT").(string)
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, path), nil
}

/*
Deprecated: This method is deprecated and will be removed soon.

Note: Please use s3_client_helper.S3Client.Delete instead.
*/
func DeleteImageByUrl(companyCode string, path string, s3Session *session.Session) (bool, error) {
	svc := s3.New(s3Session)
	bucket := viper.Get("IMAGE_UPLOAD_BUCKET").(string)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf(`esb/%s/%s`, companyCode, path)),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

/*
Deprecated: This method is deprecated and will be removed soon.
*/
func DeleteImageByKey(key string, s3Session *session.Session) (bool, error) {
	svc := s3.New(s3Session)
	bucket := viper.Get("IMAGE_UPLOAD_BUCKET").(string)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
