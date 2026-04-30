package return_ticket_repository

import (
	"core-ticket/base/helpers/query_helper/query_builder"
	"core-ticket/constants"
	"core-ticket/modules/return_ticket/return_ticket_dto"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ReturnTicketRepositoryImpl struct {
	db map[string]*sqlx.DB
}

func NewReturnTicketRepository(db map[string]*sqlx.DB) ReturnTicketRepository {
	return &ReturnTicketRepositoryImpl{db}
}

// GetReturnTicketsWithAnalysis fetches return tickets for analysis
func (r *ReturnTicketRepositoryImpl) GetReturnTicketsWithAnalysis(productTypeID int) ([]return_ticket_dto.ReturnTicket, error) {
	qb := query_builder.NewWithDB(r.db[constants.DBTicketing])
	qb.Select(
		"a.ticketNum",
		"a.companyCode",
		"a.outletCode",
		"a.outletName",
		"a.issue",
		"a.issueId",
		"a.description",
		"a.assignDeveloper",
		"a.checkNotes",
		"a.devCheckNotes",
		"a.solution",
		"a.statusId",
		"a.createdBy",
		"a.createdDate",
		"a.updatedBy",
		"a.updatedDate",
		"b.date",
		"b.ref",
		"b.refFrom",
		"b.userId",
		"c.productTypeId",
	)
	qb.From("esb_support.tr_ticket a")
	qb.LeftJoin("esb_support.tr_ticket_history b", "a.ticketNum = b.ticketNum")
	qb.LeftJoin("esb_support.tr_ticket_product_type c", "a.ticketNum = c.ticketNum")
	qb.AndWhere("=", "c.productTypeId", productTypeID)
	qb.AndWhere("LIKE", "b.description", "%return ticket%")
	qb.AndWhereRaw("a.devCheckNotes != ''")
	qb.OrderBy("a.devCheckNotes ASC, a.ticketNum ASC")

	var tickets []return_ticket_dto.ReturnTicket
	if err := qb.All(&tickets); err != nil {
		return nil, fmt.Errorf("failed to get return tickets: %w", err)
	}

	return tickets, nil
}

// GetReturnTicketsByIssueType fetches return tickets filtered by specific issue type
func (r *ReturnTicketRepositoryImpl) GetReturnTicketsByIssueType(productTypeID int, issueType string) ([]return_ticket_dto.ReturnTicket, error) {
	qb := query_builder.NewWithDB(r.db[constants.DBTicketing])
	qb.Select(
		"a.ticketNum",
		"a.companyCode",
		"a.outletCode",
		"a.outletName",
		"a.issue",
		"a.issueId",
		"a.description",
		"a.assignDeveloper",
		"a.checkNotes",
		"a.devCheckNotes",
		"a.solution",
		"a.statusId",
		"a.createdBy",
		"a.createdDate",
		"a.updatedBy",
		"a.updatedDate",
		"b.id",
		"b.date",
		"b.ref",
		"b.refFrom",
		"b.userId",
		"b.description",
		"c.id",
		"c.productTypeId",
	)
	qb.From("esb_support.tr_ticket a")
	qb.LeftJoin("esb_support.tr_ticket_history b", "a.ticketNum = b.ticketNum")
	qb.LeftJoin("esb_support.tr_ticket_product_type c", "a.ticketNum = c.ticketNum")
	qb.AndWhere("=", "c.productTypeId", productTypeID)
	qb.AndWhere("LIKE", "b.description", "%return ticket%")
	qb.AndWhere("LIKE", "a.devCheckNotes", issueType)
	qb.OrderBy("a.ticketNum ASC")

	var tickets []return_ticket_dto.ReturnTicket
	if err := qb.All(&tickets); err != nil {
		return nil, fmt.Errorf("failed to get return tickets: %w", err)
	}

	return tickets, nil
}
