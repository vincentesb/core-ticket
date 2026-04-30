package healthcheck_dto

type HealthCheckQueryRequest struct {
	ServerCode string `form:"serverCode" binding:"required"`
}
