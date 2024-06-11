package postgres

import (
	"CRM/go/subsService/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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
		log.Fatalf("error connection to database: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("error ping to database: %v", err)
	}
}
