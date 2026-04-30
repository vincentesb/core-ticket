package healthcheck

import (
	"core-ticket/base/helpers/gin_helper"
)

func Router(router *gin_helper.Router) {
	healthcheckHandler, _ := InitializeHealthCheckHandler(router.DBInstances())

	gin_helper.GET(router.Group("/health"), "", healthcheckHandler.HealthCheck)
}
