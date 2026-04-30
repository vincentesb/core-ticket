//go:build wireinject
// +build wireinject

package healthcheck

import (
	"core-ticket/modules/common/healthcheck/healthcheck_handler"
	"core-ticket/modules/common/healthcheck/healthcheck_service"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func InitializeHealthCheckHandler(dbConnections map[string]*sqlx.DB) (healthcheck_handler.HealthCheckHandler, error) {
	wire.Build(
		healthcheck_service.NewHealthCheckService,
		healthcheck_handler.NewHealthCheckHandler,
	)
	return nil, nil
}
