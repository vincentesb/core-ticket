package google_chat_helper

import (
	"bytes"
	constants "core-ticket/base/helpers/google_chat_helper/google_chat_constants"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type GoogleChatConfig struct {
	ErrorType constants.NotificationType `validate:"required"`
	Payload   map[string]interface{}     `validate:"required,dive"`

	webhookURL string `validate:"required"`
}

func SendMessage(config GoogleChatConfig) error {
	gChat := GoogleChatConfig{
		ErrorType: config.ErrorType,
		Payload:   config.Payload,

		webhookURL: viper.GetString("GOOGLE_CHAT_WEBHOOK_URL"),
	}

	validate := validator.New()
	err := validate.Struct(gChat)
	if err != nil {
		log.Printf("failed to send message to google chat. error: %v", err)
		return err
	}

	if !gChat.ErrorType.IsValid() {
		log.Printf("failed to send message to google chat. error: invalid error type")
		return nil
	}

	// Send message to Google Chat
	payload := generatePayload(gChat)
	err = sendToGoogleChatWebhook(gChat, payload)
	if err != nil {
		return err
	}

	return nil
}

func sendToGoogleChatWebhook(
	config GoogleChatConfig,
	payload map[string]interface{},
) error {

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, config.webhookURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	readResponseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("failed to send message to google chat. response: %v", string(readResponseBody))
	}

	return nil
}
