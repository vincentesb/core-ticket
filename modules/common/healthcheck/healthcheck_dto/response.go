package healthcheck_dto

type HealthCheckResponse struct {
	Status           string                 `json:"status"`
	Timestamp        string                 `json:"timestamp"`
	Databases        map[string]interface{} `json:"databases"`
	AvailableServers []string               `json:"available_servers"`
}
