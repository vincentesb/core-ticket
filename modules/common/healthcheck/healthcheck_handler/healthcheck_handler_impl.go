package healthcheck_handler

import (
	"core-ticket/base/helpers/gin_helper"
	"core-ticket/modules/common/healthcheck/healthcheck_dto"
	"core-ticket/modules/common/healthcheck/healthcheck_service"
)

type HealthCheckHandlerImpl struct {
	healthcheckService healthcheck_service.HealthCheckService
}

func NewHealthCheckHandler(healthcheckService healthcheck_service.HealthCheckService) HealthCheckHandler {
	return &HealthCheckHandlerImpl{healthcheckService}
}

func (handler *HealthCheckHandlerImpl) HealthCheck(c gin_helper.Context) (healthcheck_dto.HealthCheckResponse, error) {
	var request healthcheck_dto.HealthCheckQueryRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return healthcheck_dto.HealthCheckResponse{}, err
	}

	healthStatus, err := handler.healthcheckService.CheckHealth(request.ServerCode)
	if err != nil {
		return healthcheck_dto.HealthCheckResponse{}, err
	}

	return *healthStatus, nil
}
