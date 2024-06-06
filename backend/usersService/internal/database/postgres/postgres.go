package postgres

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func GetDB() *pgxpool.Pool {
	return pool
}

func init() {
	connectionString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", config.GetConfig().DB.Host, config.GetConfig().DB.Port, config.GetConfig().DB.User, config.GetConfig().DB.Password, config.GetConfig().DB.Name)

	var err error
	pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database connection error: %v", err))
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database ping error: %v", err))
	}
}
