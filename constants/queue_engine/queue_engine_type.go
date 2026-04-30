package queue_engine

type QueueType string

const (
	TypeReport QueueType = "report"
	TypeUpload QueueType = "upload"
)
