package email_helper

import "core-ticket/base/utility/nullable"

type MailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SenderName   string
	AuthEmail    string
	AuthPassword string
	Encryption   string
}

func (mc MailConfig) isValid() bool {
	return mc.SMTPHost != "" && mc.SMTPPort != 0 && mc.AuthEmail != "" && mc.AuthPassword != ""
}

type SendMailOption struct {
	To              []string
	Cc              []string
	Subject         string
	Message         string
	Attachments     []string
	IsHtmlBody      bool
	UseCustomConfig bool
	ReplyTo         nullable.String
	MailgunTag      nullable.String
}

type MailClient interface {
	GenerateMessageFromAsset(assetName string, data interface{}) (string, error)
	SetCustomConfig(config MailConfig)
	SendMailBulk(options ...SendMailOption) error
	SendMail(option SendMailOption) error
	SendMailTLS(option SendMailOption) error
	SendMailTlsBulk(options ...SendMailOption) error
	TestSendMail() error
}

type MailClientImpl struct {
	config        MailConfig
	defaultConfig MailConfig
}
