package models

import (
	"database/sql"
)

type User struct {
	Id         int            `json:"id"`
	IdAccount  int            `json:"id_account"`
	Name       sql.NullString `json:"name"`
	Surname    sql.NullString `json:"surname"`
	Patronymic sql.NullString `json:"patronymic"`
	Gender     int            `json:"gender"`
	DateBorn   sql.NullTime   `json:"date_born"`
}
