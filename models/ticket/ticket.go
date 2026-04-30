package ticket

import "time"

// Ticket represents a support ticket
type Ticket struct {
	TicketNum       string    `db:"ticketNum"`
	CompanyCode     string    `db:"companyCode"`
	OutletCode      string    `db:"outletCode"`
	OutletName      string    `db:"outletName"`
	Issue           string    `db:"issue"`
	IssueID         int       `db:"issueId"`
	Description     string    `db:"description"`
	AssignDeveloper string    `db:"assignDeveloper"`
	CheckNotes      string    `db:"checkNotes"`
	DevCheckNotes   string    `db:"devCheckNotes"`
	Solution        string    `db:"solution"`
	StatusID        int       `db:"statusId"`
	CreatedBy       string    `db:"createdBy"`
	CreatedDate     time.Time `db:"createdDate"`
	UpdatedBy       string    `db:"updatedBy"`
	UpdatedDate     time.Time `db:"updatedDate"`
}
