package google_chat_helper

import constants "core-ticket/base/helpers/google_chat_helper/google_chat_constants"

func generateServerError(request GoogleChatConfig, title string) map[string]interface{} {
	IconInvite := constants.IconInvite

	payloadHeader := generatePayloadHeader(&title, nil, nil, nil, nil)
	payloadFooter := generatePayloadFooter()
	payloadCard := generatePayloadCard()
	payloadWidget := []map[string]interface{}{
		generateLabel(request.Payload["transactionDate"].(string), &IconInvite, false),
		generateLabel(request.Payload["title"].(string), nil, false),
	}

	if errorMessage, ok := request.Payload["errorMessage"]; ok {
		payloadWidget = append(payloadWidget, generateLabel(errorMessage.(string), nil, true))
	}

	if cardsV2, ok := payloadCard["cardsV2"].([]map[string]interface{}); ok && len(cardsV2) > 0 {
		if card, ok := cardsV2[0]["card"].(map[string]interface{}); ok {
			card["header"] = payloadHeader
		}
	}

	if cardsV2, ok := payloadCard["cardsV2"].([]map[string]interface{}); ok && len(cardsV2) > 0 {
		if card, ok := cardsV2[0]["card"].(map[string]interface{}); ok {
			if section, ok := card["sections"].(map[string]interface{}); ok {
				section["widgets"] = append(payloadWidget, payloadFooter...)
			}
		}
	}

	return payloadCard
}
