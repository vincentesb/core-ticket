package return_ticket_dto

import (
	"core-ticket/base/utility/nullable"
	"time"
)

// ReturnTicket combines ticket with history and product type info
type ReturnTicket struct {
	TicketNum       string          `db:"ticketNum"`
	CompanyCode     string          `db:"companyCode"`
	OutletCode      string          `db:"outletCode"`
	OutletName      string          `db:"outletName"`
	Issue           nullable.String `db:"issue"`
	IssueID         nullable.Int    `db:"issueId"`
	Description     string          `db:"description"`
	AssignDeveloper string          `db:"assignDeveloper"`
	CheckNotes      string          `db:"checkNotes"`
	DevCheckNotes   string          `db:"devCheckNotes"`
	Solution        nullable.String `db:"solution"`
	StatusID        int             `db:"statusId"`
	CreatedBy       string          `db:"createdBy"`
	CreatedDate     time.Time       `db:"createdDate"`
	UpdatedBy       string          `db:"updatedBy"`
	UpdatedDate     time.Time       `db:"updatedDate"`
	Date            time.Time       `db:"date"`
	Ref             nullable.String `db:"ref"`
	RefFrom         nullable.String `db:"refFrom"`
	UserID          string          `db:"userId"`
	ProductTypeID   int             `db:"productTypeId"`
}

// IssueAnalysis groups similar issues by devCheckNotes
type IssueAnalysis struct {
	IssueType   string         // The devCheckNotes value (main problem classification)
	Count       int            // How many tickets have this issue type
	Tickets     []ReturnTicket // All tickets with this issue type
	MainProblem string         // Summary of the main problem based on descriptions
}

// AnalysisResult contains the full analysis of return tickets
type AnalysisResult struct {
	TotalTickets    int
	IssueGroups     []IssueAnalysis
	MostCommonIssue *IssueAnalysis // The issue type with the most occurrences
}
