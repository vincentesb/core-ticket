package email_helper

import (
	"bytes"
	"core-ticket/constants"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func NewMailClient() MailClient {
	port, _ := strconv.Atoi(viper.Get("SMTP_PORT").(string))
	return &MailClientImpl{
		defaultConfig: MailConfig{
			SMTPHost:     viper.Get("SMTP_HOST").(string),
			SMTPPort:     port,
			SenderName:   viper.Get("SMTP_USERNAME").(string),
			AuthEmail:    viper.Get("SMTP_USERNAME").(string),
			AuthPassword: viper.Get("SMTP_PASSWORD").(string),
			Encryption:   viper.Get("SMTP_ENCRYPTION").(string),
		},
	}
}

func (m *MailClientImpl) SetCustomConfig(config MailConfig) {
	if config.SMTPPort != 0 && config.SMTPHost != "" && config.AuthEmail != "" && config.AuthPassword != "" {
		m.config = config
	} else {
		m.config = m.defaultConfig
	}
}

func (m *MailClientImpl) SendMailBulk(options ...SendMailOption) error {
	for _, option := range options {
		err := m.SendMail(option)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MailClientImpl) constructBody(option SendMailOption, config MailConfig) (string, error) {
	boundary := "boundary_" + uuid.New().String()

	headers := make(map[string]string)
	if config.SenderName != "" {
		headers["From"] = fmt.Sprintf("%s <%s>", "no-reply", config.SenderName)
	} else {
		headers["From"] = config.SenderName
	}
	headers["To"] = strings.Join(option.To, ",")
	if len(option.Cc) > 0 {
		headers["Cc"] = strings.Join(option.Cc, ",")
	}
	headers["Subject"] = option.Subject
	if option.ReplyTo.ValueOrZero() != "" {
		headers["Reply-To"] = option.ReplyTo.String
	}
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)
	if option.MailgunTag.ValueOrZero() != "" {
		headers["X-Mailgun-Tag"] = option.MailgunTag.String
	} else {
		headers["X-Mailgun-Tag"] = constants.DefaultMailgunTag
	}

	var emailBody strings.Builder
	for k, v := range headers {
		emailBody.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	emailBody.WriteString("\r\n")

	emailBody.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	if option.IsHtmlBody {
		emailBody.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	} else {
		emailBody.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	}
	emailBody.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	emailBody.WriteString(option.Message)
	emailBody.WriteString("\r\n")

	for _, attachment := range option.Attachments {
		fileData, err := os.ReadFile(attachment)
		if err != nil {
			return "", err
		}

		emailBody.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		emailBody.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filepath.Base(attachment)))
		emailBody.WriteString("Content-Type: application/octet-stream\r\n")
		emailBody.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		emailBody.WriteString(base64.StdEncoding.EncodeToString(fileData))
		emailBody.WriteString("\r\n")
	}

	// End boundary
	emailBody.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return emailBody.String(), nil
}

func (m *MailClientImpl) SendMail(option SendMailOption) error {
	config := m.defaultConfig

	if option.UseCustomConfig && m.config.isValid() {
		config = m.config
	}

	if !config.isValid() {
		return fmt.Errorf("invalid SMTP configuration")
	}

	validEncryptionTLS := []string{"TLS"}
	if slices.Contains(validEncryptionTLS, strings.ToUpper(config.Encryption)) {
		return m.SendMailTLS(option)
	}

	validEncryptionSSL := []string{"SSL"}

	if slices.Contains(validEncryptionSSL, strings.ToUpper(config.Encryption)) {
		return m.SendMailSSL(option)
	}

	auth := smtp.PlainAuth("", config.AuthEmail, config.AuthPassword, config.SMTPHost)
	smtpAddr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	body, err := m.constructBody(option, config)
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpAddr, auth, config.AuthEmail, append(option.To, option.Cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}

func (m *MailClientImpl) TestSendMail() error {
	option := SendMailOption{
		To:              []string{m.config.AuthEmail},
		Subject:         "Email Sending Test",
		Message:         "Dear User,<br><br><p>If you received this email, your email credential configuration is successfully set in ESB Core.</p><p>(This is an auto generated Email)</p><br><br>Thank You.",
		UseCustomConfig: true,
		IsHtmlBody:      true,
	}

	return m.SendMail(option)
}

func (m *MailClientImpl) SendMailTLS(option SendMailOption) error {
	config := m.defaultConfig

	if option.UseCustomConfig && m.config.isValid() {
		config = m.config
	}

	if !config.isValid() {
		return fmt.Errorf("invalid SMTP configuration")
	}

	auth := smtp.PlainAuth("", config.AuthEmail, config.AuthPassword, config.SMTPHost)
	smtpAddr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	// Connect to the SMTP server with TLS encryption
	conn, err := smtp.Dial(smtpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Start TLS encryption
	tlsCfg := &tls.Config{ServerName: config.SMTPHost}
	if err := conn.StartTLS(tlsCfg); err != nil {
		return err
	}

	// Authenticate with the server
	if err = conn.Auth(auth); err != nil {
		return err
	}

	// Send the email
	if err = conn.Mail(config.AuthEmail); err != nil {
		return err
	}
	for _, recipient := range option.To {
		if err = conn.Rcpt(recipient); err != nil {
			return err
		}
	}
	wc, err := conn.Data()
	if err != nil {
		return err
	}

	body, err := m.constructBody(option, config)
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(body))
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *MailClientImpl) SendMailTlsBulk(options ...SendMailOption) error {
	for _, option := range options {
		err := m.SendMailTLS(option)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MailClientImpl) SendMailSSL(option SendMailOption) error {
	config := m.defaultConfig

	if option.UseCustomConfig && m.config.isValid() {
		config = m.config
	}

	if !config.isValid() {
		return fmt.Errorf("invalid SMTP configuration")
	}

	auth := smtp.PlainAuth("", config.AuthEmail, config.AuthPassword, config.SMTPHost)
	smtpAddr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	tlsCfg := &tls.Config{ServerName: config.SMTPHost}
	conn, err := tls.Dial("tcp", smtpAddr, tlsCfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Quit()

	// Authenticate with the server
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Send the email
	if err = client.Mail(config.AuthEmail); err != nil {
		return err
	}
	for _, recipient := range option.To {
		if err = client.Rcpt(recipient); err != nil {
			return err
		}
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}

	body, err := m.constructBody(option, config)
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(body))
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *MailClientImpl) GenerateMessageFromAsset(assetName string, data interface{}) (string, error) {
	templateData, err := os.ReadFile(fmt.Sprintf("assets/email_template/%s.html", assetName))
	if err != nil {
		return "", err
	}

	t, err := template.New("emailTemplate").Parse(string(templateData))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
