package database

import (
	"core-ticket/config"
	"core-ticket/constants"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type ClientDB struct {
	ServerCode string `db:"serverCode"`
	Hostname   string `db:"hostName"`
	Username   string `db:"username"`
	Password   string `db:"password"`
}

func NewDB(cfg config.AppConfig) map[string]*sqlx.DB {
	dbConnections := make(map[string]*sqlx.DB)
	dbMain := newMainDB(cfg)
	dbTicketing := newTicketingDB(cfg)

	dbConnections[constants.DBMain] = dbMain
	dbConnections[constants.DBTicketing] = dbTicketing
	return dbConnections
}

func newTicketingDB(cfg config.AppConfig) *sqlx.DB {
	dbHost := cfg.TicketingDbHost
	dbUsername := cfg.TicketingDbUsername
	dbPassword := cfg.TicketingDbPassword
	dbName := cfg.TicketingDbName
	timezone := cfg.Timezone

	db := connectDB(dbHost, dbName, dbUsername, dbPassword, timezone)
	return db
}

func newMainDB(cfg config.AppConfig) *sqlx.DB {
	dbHost := cfg.MainDbHost
	dbUsername := cfg.MainDbUsername
	dbPassword := cfg.MainDbPassword
	dbName := cfg.MainDbName
	timezone := cfg.Timezone

	db := connectDB(dbHost, dbName, dbUsername, dbPassword, timezone)
	return db
}

func connectDB(dbHost string, dbName string, dbUsername string, dbPassword string, timezone string) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true&loc=%s&multiStatements=true",
		dbUsername, dbPassword, dbHost, dbName, url.QueryEscape(timezone))
	fmt.Println("dsn", dsn)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database %s %v", dbName, err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database %s %v", dbName, err))
	}
	fmt.Printf("Database %s successfully connected\n", dbName)

	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}
