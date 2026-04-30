package healthcheck_service

import (
	"core-ticket/base/helpers/error_helper"
	"core-ticket/constants"
	"core-ticket/constants/error_code"
	"core-ticket/modules/common/healthcheck/healthcheck_dto"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBStatus string

const (
	DBStatusHealthy   DBStatus = "healthy"
	DBStatusUnhealthy DBStatus = "unhealthy"
	DBStatusNotFound  DBStatus = "not_found"
)

type ClientDB struct {
	ServerCode string `db:"serverCode"`
	Hostname   string `db:"hostName"`
	Username   string `db:"username"`
	Password   string `db:"password"`
}

type HealthCheckServiceImpl struct {
	dbConnections map[string]*sqlx.DB
}

func NewHealthCheckService(dbConnections map[string]*sqlx.DB) HealthCheckService {
	return &HealthCheckServiceImpl{
		dbConnections: dbConnections,
	}
}

func (service *HealthCheckServiceImpl) CheckHealth(serverCode string) (*healthcheck_dto.HealthCheckResponse, error) {
	dbHealth := make(map[string]interface{})
	databaseStatus := DBStatusHealthy

	if db, exists := service.dbConnections[serverCode]; exists {
		if err := db.Ping(); err != nil {
			dbHealth[serverCode] = map[string]interface{}{
				"status": DBStatusUnhealthy,
				"error":  err.Error(),
			}
			databaseStatus = DBStatusUnhealthy
		} else {
			dbHealth[serverCode] = map[string]interface{}{
				"status": DBStatusHealthy,
			}
		}
	} else {
		dbHealth[serverCode] = map[string]interface{}{
			"status": DBStatusNotFound,
			"error":  fmt.Sprintf("Server code '%s' not found in active connections", serverCode),
		}
		databaseStatus = DBStatusNotFound
	}

	availableServers, err := service.getAvailableServerCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve available server codes: %v", err)
	}

	healthStatus := healthcheck_dto.HealthCheckResponse{
		Status:           string(databaseStatus),
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		Databases:        dbHealth,
		AvailableServers: availableServers,
	}

	if databaseStatus != DBStatusHealthy {
		return nil, error_helper.New(errors.New("database is not healthy"), error_code.ServiceUnavailable)
	}

	return &healthStatus, nil
}

func (service *HealthCheckServiceImpl) getAvailableServerCodes() ([]string, error) {
	mainDB, exists := service.dbConnections[constants.DBMain]
	if !exists {
		return nil, fmt.Errorf("main database connection not found")
	}

	ticketingDB, exists := service.dbConnections[constants.DBTicketing]
	if !exists {
		return nil, fmt.Errorf("ticketing database connection not found")
	}

	serverCodes := []string{constants.DBMain, constants.DBTicketing}
	if mainDB != ticketingDB {
		serverCodes = append(serverCodes, constants.ESBFnbDB)
	}

	return serverCodes, nil
}
