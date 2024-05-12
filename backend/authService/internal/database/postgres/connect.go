package postgres

import (
	"CRM/go/authService/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	cfg := config.GetConfig()
	connectionString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to postgres!")
	return db, nil
}
