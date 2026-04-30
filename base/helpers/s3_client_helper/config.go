package s3_client_helper

import (
	"github.com/spf13/viper"
)

type S3Config struct {
	EndPoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
}

/*
getConfig retrieves S3 configuration values from the environment using viper and returns a pointer to an S3Config struct containing the retrieved values.

Returns:

	*S3Config: A pointer to an S3Config struct with the retrieved S3 configuration values.
*/
func getConfig() *S3Config {
	return &S3Config{
		EndPoint:  viper.GetString("S3_ENDPOINT"),
		Bucket:    viper.GetString("S3_BUCKET"),
		AccessKey: viper.GetString("S3_ACCESS_KEY"),
		SecretKey: viper.GetString("S3_SECRET_KEY"),
		Region:    viper.GetString("S3_REGION"),
	}
}

// getReportConfig returns a configuration for report uploads. If REPORT_S3_ENDPOINT is set it will use REPORT_S3_* env vars.
// Returns nil if REPORT_S3_ENDPOINT is not set (so caller can fallback to default config).
func getReportConfig() *S3Config {
	if endpoint := viper.GetString("REPORT_S3_ENDPOINT"); endpoint != "" {
		return &S3Config{
			EndPoint:  endpoint,
			Bucket:    viper.GetString("REPORT_S3_BUCKET"),
			AccessKey: viper.GetString("REPORT_S3_ACCESS_KEY"),
			SecretKey: viper.GetString("REPORT_S3_SECRET_KEY"),
			Region:    viper.GetString("REPORT_S3_REGION"),
		}
	}
	return nil
}
