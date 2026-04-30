package email_helper_v2

import (
	"context"
	"core-ticket/constants"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/spf13/viper"
)

func NewMailClientV2() MailClientV2 {
	return &MailClientV2Impl{
		config: MailConfigV2{
			MailgunApiKey: viper.GetString("MAILGUN_API_KEY"),
			MailgunDomain: viper.GetString("MAILGUN_DOMAIN"),
		},
	}
}

func (m *MailClientV2Impl) SendMailBulk(options ...SendMailOptionV2) ([]string, error) {
	msgIds := make([]string, len(options))
	for _, option := range options {
		msgId, err := m.SendMail(option)
		if err != nil {
			return nil, err
		}
		msgIds = append(msgIds, msgId)
	}

	return msgIds, nil
}

func (m *MailClientV2Impl) SendMail(option SendMailOptionV2) (string, error) {
	if !m.config.isValid() {
		return "", fmt.Errorf("invalid mailgun config")
	}

	mg := mailgun.NewMailgun(m.config.MailgunDomain, m.config.MailgunApiKey)

	from := fmt.Sprintf("no-reply@%s", m.config.MailgunDomain)
	var message *mailgun.Message
	if option.IsHtmlBody {
		message = mailgun.NewMessage(
			from,
			option.Subject,
			"",
			option.To...,
		)
	} else {
		message = mailgun.NewMessage(
			from,
			option.Subject,
			option.Message,
			option.To...,
		)
	}

	// CC
	if len(option.Cc) > 0 {
		for _, mail := range option.Cc {
			message.AddCC(mail)
		}
	}

	// Set HTML content
	if option.IsHtmlBody {
		message.SetHTML(option.Message)
	}

	// ReplyTo
	if option.ReplyTo.Valid {
		message.SetReplyTo(option.ReplyTo.String)
	}

	// Tagging
	if option.MailgunTag.Valid {
		message.AddTag(option.MailgunTag.String)
	} else {
		message.AddTag(constants.DefaultMailgunTag)
	}

	// Attachments
	for _, file := range option.Attachments {
		message.AddAttachment(file)
	}

	// Send with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, msgId, err := mg.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("failed to send mail: %w", err)
	}

	return msgId, nil
}
