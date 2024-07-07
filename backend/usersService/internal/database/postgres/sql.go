package postgres

import (
	"CRM/go/usersService/internal/proto/usersService"
	"CRM/go/usersService/pkg/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Registration(user *usersService.User) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "INSERT INTO users (id_account) VALUES ($1)", user.Id)
}

func GetUser(user *usersService.User) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT users.name, users.surname, users.patronymic, users.gender, users.date_born, users.position, CASE WHEN users.position = 'trainer' THEN STRING_AGG(CONCAT(trainers.exp, ';', trainers.sport, ';', trainers.achievement), '|') ELSE NULL END AS trainer_info FROM users LEFT JOIN trainers ON users.id_account = trainers.id_account AND users.position = 'trainer' WHERE users.id_account = $1 GROUP BY users.name, users.surname, users.patronymic, users.gender, users.date_born, users.position", user.Id)
}

func GetTrainers() (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT users.id_account, users.name, users.surname, users.patronymic, users.gender, users.date_born, users.position, STRING_AGG(CONCAT(trainers.exp, ';', trainers.sport, ';', trainers.achievement), '|') AS trainer_info FROM users JOIN trainers ON users.id_account = trainers.id_account WHERE users.position = 'trainer' GROUP BY users.id, users.name, users.surname, users.patronymic, users.gender, users.date_born, users.position")
}

func UpdateUser(user *usersService.User) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE users SET name = $1, surname = $2, patronymic = $3, gender = $4, date_born = $5 WHERE id_account = $6", user.Name, user.Surname, user.Patronymic, user.Gender, utils.ConvertTimestampToSQLNullTime(user.DateBorn), user.Id)
}
