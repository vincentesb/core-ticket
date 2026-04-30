package queue_engine

const (
	EnvKeyReportApiUrl   = "QUEUE_API_URL"
	EnvKeyReportApiToken = "QUEUE_API_TOKEN"
	EnvKeyUploadApiUrl   = "UPLOAD_QUEUE_API_URL"
	EnvKeyUploadApiToken = "UPLOAD_QUEUE_API_TOKEN"
)

const MaxUploadQueueCount = 10

const (
	SourceCore        = "core"
	SourceCoreService = "core-service"
)
