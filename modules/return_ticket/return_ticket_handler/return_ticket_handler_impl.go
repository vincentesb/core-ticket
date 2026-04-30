package return_ticket_handler

import (
	"core-ticket/base/helpers/gin_helper"
	"core-ticket/modules/return_ticket/return_ticket_dto"
	"core-ticket/modules/return_ticket/return_ticket_service"
)

type ReturnTicketHandlerImpl struct {
	returnTicketService return_ticket_service.ReturnTicketService
}

func NewReturnTicketHandler(
	returnTicketService return_ticket_service.ReturnTicketService,
) ReturnTicketHandler {
	return &ReturnTicketHandlerImpl{
		returnTicketService,
	}
}

func (h *ReturnTicketHandlerImpl) GetReturnTicketsWithAnalysis(c gin_helper.Context) (*return_ticket_dto.AnalysisResult, error) {
	var request return_ticket_dto.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, err
	}
	return h.returnTicketService.GetReturnTicketsWithAnalysis(request.ProductTypeID)
}

func (h *ReturnTicketHandlerImpl) GetReturnTicketsByIssueType(c gin_helper.Context) ([]return_ticket_dto.ReturnTicket, error) {
	var request return_ticket_dto.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, err
	}
	return h.returnTicketService.GetReturnTicketsByIssueType(request.ProductTypeID, request.IssueType)
}
