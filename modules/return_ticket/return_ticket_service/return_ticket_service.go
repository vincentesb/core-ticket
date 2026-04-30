package return_ticket_service

import "core-ticket/modules/return_ticket/return_ticket_dto"

type ReturnTicketService interface {
	GetReturnTicketsWithAnalysis(productTypeID int) (*return_ticket_dto.AnalysisResult, error)
	GetReturnTicketsByIssueType(productTypeID int, issueType string) ([]return_ticket_dto.ReturnTicket, error)
}
