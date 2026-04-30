package email_helper_v2

import (
	"core-ticket/base/utility/nullable"
)

type MailConfigV2 struct {
	MailgunApiKey string
	MailgunDomain string
}

func (mc MailConfigV2) isValid() bool {
	return mc.MailgunApiKey != "" && mc.MailgunDomain != ""
}

type SendMailOptionV2 struct {
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

type MailClientV2 interface {
	SendMail(option SendMailOptionV2) (string, error)
	SendMailBulk(options ...SendMailOptionV2) ([]string, error)
}

type MailClientV2Impl struct {
	config MailConfigV2
}
