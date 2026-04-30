package ticket_history

import "time"

// TicketHistory represents the history of a ticket
type TicketHistory struct {
	ID          int       `db:"id"`
	TicketNum   string    `db:"ticketNum"`
	Date        time.Time `db:"date"`
	Ref         string    `db:"ref"`
	RefFrom     string    `db:"refFrom"`
	UserID      string    `db:"userId"`
	Description string    `db:"description"`
}
