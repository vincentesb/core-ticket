package healthcheck_handler

import (
	"core-ticket/base/helpers/gin_helper"
	"core-ticket/modules/common/healthcheck/healthcheck_dto"
)

type HealthCheckHandler interface {
	HealthCheck(c gin_helper.Context) (healthcheck_dto.HealthCheckResponse, error)
}
