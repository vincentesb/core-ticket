package return_ticket

import "core-ticket/base/helpers/gin_helper"

func Router(router *gin_helper.Router) {
	handler := InitializeReturnTicketHandler(router.DBInstances())

	rg := router.Group("/return-ticket")
	{
		gin_helper.GET(rg, "", handler.GetReturnTicketsByIssueType)
		gin_helper.GET(rg, "/analyze", handler.GetReturnTicketsWithAnalysis)
	}
}
