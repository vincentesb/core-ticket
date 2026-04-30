package ticket_product_type

// TicketProductType represents product type associated with a ticket
type TicketProductType struct {
	ID            int    `db:"id"`
	TicketNum     string `db:"ticketNum"`
	ProductTypeID int    `db:"productTypeId"`
}
