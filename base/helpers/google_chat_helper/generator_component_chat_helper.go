package google_chat_helper

import constants "core-ticket/base/helpers/google_chat_helper/google_chat_constants"

func generatePayload(request GoogleChatConfig) map[string]interface{} {
	var payload map[string]interface{}

	titleText := request.ErrorType.String()
	switch request.ErrorType {
	case constants.TypeInternalServerError:
		payload = generateServerError(request, titleText)
	case constants.TypeFailedConsumeDisbursementData:
		payload = generateFailedConsumeDisbursementData(request, titleText)
	}

	return payload
}

func generatePayloadCard() map[string]interface{} {
	return map[string]interface{}{
		"cardsV2": []map[string]interface{}{
			{
				"cardId": "unique-card-id",
				"card": map[string]interface{}{
					"header": map[string]interface{}{},
					"sections": map[string]interface{}{
						"widgets": []map[string]interface{}{},
					},
				},
			},
		},
	}
}

func generatePayloadHeader(
	title *string,
	subtitle *string,
	imageUrl *string,
	imageType *string,
	imageAltText *string,
) map[string]interface{} {
	if title == nil {
		value := constants.DefaultTitleHeader
		title = &value
	}
	if subtitle == nil {
		value := constants.DefaultSubtitleHeader
		subtitle = &value
	}
	if imageUrl == nil {
		value := constants.DefaultImageUrlHeader
		imageUrl = &value
	}
	if imageType == nil {
		value := constants.DefaultImageTypeHeader
		imageType = &value
	}
	if imageAltText == nil {
		value := constants.DefaultImageAltTextHeader
		imageAltText = &value
	}

	return map[string]interface{}{
		"title":        title,
		"subtitle":     subtitle,
		"imageUrl":     imageUrl,
		"imageType":    imageType,
		"imageAltText": imageAltText,
	}
}

func generatePayloadFooter() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"divider": map[string]interface{}{},
		},
		generateLabel("👀 (eye) emoji to assign an issue to yourself.", nil, false),
		generateLabel("✅ (check) if the issue has been resolved.", nil, false),
		{
			"divider": map[string]interface{}{},
		},
		generateLabel("Please reply to this thread if you want to discuss.", nil, false),
	}
}

func generateLabel(label string, icon *string, wrapText bool) map[string]interface{} {
	payload := map[string]interface{}{
		"decoratedText": map[string]interface{}{
			"startIcon": map[string]interface{}{
				"knownIcon": icon,
			},
			"text":     label,
			"wrapText": wrapText,
		},
	}
	if icon == nil {
		if decoratedText, ok := payload["decoratedText"].(map[string]interface{}); ok {
			delete(decoratedText, "startIcon")
		}
	}

	return payload
}

func generateDivider() map[string]interface{} {
	return map[string]interface{}{
		"divider": map[string]interface{}{},
	}
}

func generateParagraph(text string) map[string]interface{} {
	return map[string]interface{}{
		"textParagraph": map[string]interface{}{
			"text": text,
		},
	}
}
