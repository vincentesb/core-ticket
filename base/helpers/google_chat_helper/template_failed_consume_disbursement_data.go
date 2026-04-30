package google_chat_helper

import constants "core-ticket/base/helpers/google_chat_helper/google_chat_constants"

func generateFailedConsumeDisbursementData(request GoogleChatConfig, title string) map[string]interface{} {
	IconInvite := constants.IconInvite
	IconStore := constants.IconStore
	IconTicket := constants.IconTicket
	IconStar := constants.IconStar
	IconDescription := constants.IconDescription

	payloadHeader := generatePayloadHeader(&title, nil, nil, nil, nil)
	payloadFooter := generatePayloadFooter()
	payloadCard := generatePayloadCard()
	payloadWidget := []map[string]interface{}{
		generateLabel(request.Payload["transactionDate"].(string), &IconInvite, false),
		generateLabel(request.Payload["companyCode"].(string), &IconStore, false),
	}

	if disburseNum, ok := request.Payload["disburseNum"]; ok {
		payloadWidget = append(payloadWidget, generateLabel(disburseNum.(string), &IconTicket, false))
	}

	payloadWidget = append(payloadWidget, generateDivider())

	if titleText, ok := request.Payload["title"]; ok {
		titleStr := titleText.(string)
		if len(titleStr) <= 40 {
			payloadWidget = append(payloadWidget, generateLabel(titleStr, &IconStar, false))
		} else {
			payloadWidget = append(payloadWidget, generateParagraph(titleStr))
		}
	}

	if descriptionText, ok := request.Payload["description"]; ok {
		description := descriptionText.(string)
		if len(description) <= 40 {
			payloadWidget = append(payloadWidget, generateLabel(description, &IconDescription, false))
		} else {
			payloadWidget = append(payloadWidget, generateParagraph(description))
		}
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
