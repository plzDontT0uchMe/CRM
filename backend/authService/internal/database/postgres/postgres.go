package postgres

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func init() {
	connectionString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", config.GetConfig().DB.DBHost, config.GetConfig().DB.DBPort, config.GetConfig().DB.DBUser, config.GetConfig().DB.DBPassword, config.GetConfig().DB.DBName)
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database connection error: %v", err))
		return
	}
	err = db.Ping()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database ping error: %v", err))
		return
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Second * 60)
	db.SetConnMaxIdleTime(time.Second * 60)
}
