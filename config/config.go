package config

import (
	"github.com/spf13/viper"
)

// AppConfig adalah struktur untuk menyimpan konfigurasi aplikasi.
type AppConfig struct {
	AppEnv              string `mapstructure:"APP_ENV"`
	AppFrontendUrl      string `mapstructure:"APP_FRONTEND_URL"`
	AppHost             string `mapstructure:"APP_HOST"`
	AppPort             string `mapstructure:"APP_PORT"`
	AppPrivateHost      string `mapstructure:"APP_PRIVATE_HOST"`
	AppJwtSecret        string `mapstructure:"APP_JWT_SECRET"`
	AppJwtTokenLifeSpan string `mapstructure:"APP_JWT_TOKEN_LIFE_SPAN"`
	MainDbHost          string `mapstructure:"MAIN_DB_HOST"`
	MainDbUsername      string `mapstructure:"MAIN_DB_USERNAME"`
	MainDbPassword      string `mapstructure:"MAIN_DB_PASSWORD"`
	MainDbName          string `mapstructure:"MAIN_DB_NAME"`
	DbPort              string `mapstructure:"DB_PORT"`
	DbName              string `mapstructure:"DB_NAME"`
	SentryDsn           string `mapstructure:"SENTRY_DSN"`
	SentrySampleRate    string `mapstructure:"SENTRY_SAMPLE_RATE"`
	SentryProfileRate   string `mapstructure:"SENTRY_PROFILE_RATE"`
	SecurityKey         string `mapstructure:"SECURITY_KEY"`
	SmtpHost            string `mapstructure:"SMTP_HOST"`
	SmtpPort            string `mapstructure:"SMTP_PORT"`
	SmtpEncryption      string `mapstructure:"SMTP_ENCRYPTION"`
	SmtpUsername        string `mapstructure:"SMTP_USERNAME"`
	SmtpPassword        string `mapstructure:"SMTP_PASSWORD"`
	TicketingDbHost     string `mapstructure:"TICKETING_DB_HOST"`
	TicketingDbUsername string `mapstructure:"TICKETING_DB_USERNAME"`
	TicketingDbPassword string `mapstructure:"TICKETING_DB_PASSWORD"`
	TicketingDbName     string `mapstructure:"TICKETING_DB_NAME"`

	Timezone            string `mapstructure:"TZ"`
	GoAppApiUrl         string `mapstructure:"GO_APP_API_URL:"`
	QueueApiUrl         string `mapstructure:"QUEUE_API_URL"`
	QueueApiToken       string `mapstructure:"QUEUE_API_TOKEN"`
	UploadQueueApiUrl   string `mapstructure:"UPLOAD_QUEUE_API_URL"`
	UploadQueueApiToken string `mapstructure:"UPLOAD_QUEUE_API_TOKEN"`

	BasicAuthUsername string `mapstructure:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword string `mapstructure:"BASIC_AUTH_PASSWORD"`

	MyEsbUrl          string `mapstructure:"MY_ESB_URL"`
	MyEsbRestUsername string `mapstructure:"MY_ESB_REST_USERNAME"`
	MyEsbRestPassword string `mapstructure:"MY_ESB_REST_PASSWORD"`

	MailgunApiKey string `mapstructure:"MAILGUN_API_KEY"`
	MailgunDomain string `mapstructure:"MAILGUN_DOMAIN"`

	EsbGoodsUrl string `mapstructure:"ESB_GOODS_URL"`

	AppRefreshSecret string `mapstructure:"APP_REFRESH_SECRET"`
}

func LoadConfig(cfg *AppConfig) error {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("APP_FRONTEND_URL", "https://dev5.esb.co.id/erp-refactor")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "3003")
	viper.SetDefault("TZ", "Asia/Jakarta")

	viper.SetDefault("MAIN_DB_HOST", "localhost")
	viper.SetDefault("MAIN_DB_NAME", "esb_main")

	viper.SetDefault("SENTRY_DSN", "")
	viper.SetDefault("SENTRY_SAMPLE_RATE", "1.0")
	viper.SetDefault("SENTRY_PROFILE_RATE", "1.0")

	viper.SetDefault("SMTP_HOST", "mail.mailer-esb.com")
	viper.SetDefault("SMTP_PORT", "465")
	viper.SetDefault("SMTP_ENCRYPTION", "ssl")
	viper.SetDefault("SMTP_USERNAME", "")
	viper.SetDefault("SMTP_PASSWORD", "")

	viper.SetDefault("APP_JWT_TOKEN_LIFE_SPAN", 1)

	viper.SetDefault("SECURITY_KEY", "")

	viper.SetDefault("IMAGE_UPLOAD_ENDPOINT", "")
	viper.SetDefault("IMAGE_UPLOAD_KEY", "")
	viper.SetDefault("IMAGE_UPLOAD_SECRET", "")
	viper.SetDefault("IMAGE_UPLOAD_BUCKET", "")
	viper.SetDefault("IMAGE_UPLOAD_REGION", "sgp1")

	viper.SetDefault("S3_ENDPOINT", "https://oss-ap-southeast-5.aliyuncs.com")
	viper.SetDefault("S3_BUCKET", "esb-bucket-dev")
	viper.SetDefault("S3_ACCESS_KEY", "")
	viper.SetDefault("S3_SECRET_KEY", "")
	viper.SetDefault("S3_REGION", "sgp1")

	viper.SetDefault("GO_APP_API_URL", "")

	viper.SetDefault("SHARED_SERVICE_BASE_URL", "")

	viper.SetDefault("OY_PAYMENT_GATEWAY_MIN_LIMIT", 10000)

	viper.SetDefault("CRM_S3_ENDPOINT", "https://oss-ap-southeast-5.aliyuncs.com")
	viper.SetDefault("CRM_S3_BUCKET", "esb-bucket-dev")
	viper.SetDefault("CRM_S3_ACCESS_KEY", "")
	viper.SetDefault("CRM_S3_SECRET_KEY", "")
	viper.SetDefault("CRM_S3_REGION", "sgp1")
	viper.SetDefault("CRM_S3_ARCHIVE_DIR", "archive")

	viper.SetDefault("MY_ESB_URL", "")
	viper.SetDefault("MY_ESB_REST_USERNAME", "")
	viper.SetDefault("MY_ESB_REST_PASSWORD", "")

	viper.SetDefault("MAILGUN_API_KEY", "")
	viper.SetDefault("MAILGUN_DOMAIN", "")

	viper.SetDefault("ESB_GOODS_URL", "https://dev7.esb.co.id/esb-goods/frontend/web")

	viper.SetConfigName("app.local")
	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigName("app")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	return viper.Unmarshal(cfg)
}
