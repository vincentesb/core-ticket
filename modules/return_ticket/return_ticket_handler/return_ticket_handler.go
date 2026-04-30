package return_ticket_handler

import (
	"core-ticket/base/helpers/gin_helper"
	"core-ticket/modules/return_ticket/return_ticket_dto"
)

type ReturnTicketHandler interface {
	GetReturnTicketsWithAnalysis(c gin_helper.Context) (*return_ticket_dto.AnalysisResult, error)
	GetReturnTicketsByIssueType(c gin_helper.Context) ([]return_ticket_dto.ReturnTicket, error)
}
