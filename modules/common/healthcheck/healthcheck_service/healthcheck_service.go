package healthcheck_service

import "core-ticket/modules/common/healthcheck/healthcheck_dto"

type HealthCheckService interface {
	CheckHealth(serverCode string) (*healthcheck_dto.HealthCheckResponse, error)
}
