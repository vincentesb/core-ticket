//go:build wireinject
// +build wireinject

// wire:wireinject
package return_ticket

import (
	"core-ticket/modules/return_ticket/return_ticket_handler"
	"core-ticket/modules/return_ticket/return_ticket_repository"
	"core-ticket/modules/return_ticket/return_ticket_service"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func InitializeReturnTicketService(db map[string]*sqlx.DB) return_ticket_service.ReturnTicketService {
	wire.Build(
		return_ticket_repository.NewReturnTicketRepository,
		return_ticket_service.NewReturnTicketService,
	)
	return nil
}

func InitializeReturnTicketHandler(db map[string]*sqlx.DB) return_ticket_handler.ReturnTicketHandler {
	wire.Build(
		InitializeReturnTicketService,

		return_ticket_handler.NewReturnTicketHandler,
	)
	return nil
}
