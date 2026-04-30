package database_constants

import "core-ticket/config"

var (
	ESB_MAIN = "esb_main"
	ESB_FNB  = "esb_fnb"
)

const (
	TRANSACTION_TIME_OUT    = 60  // used to set context time out in seconds
	SYNC_DATA_TIME_OUT      = 10  // used to set context time out in seconds when sync data
	RUN_ASYNC_DATA_TIME_OUT = 600 // used to set context time out in seconds

)

const MAX_CHUNK_SIZE = 500 // Maximal total data when insert batch in one batch.

func InitDatabaseName(cfg config.AppConfig) {
	if cfg.MainDbName != "" {
		ESB_MAIN = cfg.MainDbName
	}
}
