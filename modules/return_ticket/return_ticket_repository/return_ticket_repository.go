package return_ticket_repository

import (
	"core-ticket/modules/return_ticket/return_ticket_dto"
)

type ReturnTicketRepository interface {
	GetReturnTicketsWithAnalysis(productTypeID int) ([]return_ticket_dto.ReturnTicket, error)
	GetReturnTicketsByIssueType(productTypeID int, issueType string) ([]return_ticket_dto.ReturnTicket, error)
}
