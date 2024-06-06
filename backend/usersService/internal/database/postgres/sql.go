package postgres

import (
	"CRM/go/usersService/internal/logger"
	"CRM/go/usersService/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func CreateUser(user *models.User) (error, int) {
	_, err := GetDB().Exec(context.Background(), "INSERT INTO users (id_account, gender) VALUES ($1, $2)", user.IdAccount, user.Gender)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err), "UserIdAccount", user.IdAccount)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func GetUser(user *models.User) (error, int) {
	err := GetDB().QueryRow(context.Background(), "SELECT name, surname, patronymic, gender, date_born FROM users WHERE id_account = $1", user.IdAccount).Scan(&user.Name, &user.Surname, &user.Patronymic, &user.Gender, &user.DateBorn)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.CreateLog("error", fmt.Sprintf("user not found: %v", err), "UserIdAccount", user.IdAccount)
			return err, http.StatusBadRequest
		}
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err), "UserIdAccount", user.IdAccount)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func UpdateUser(user models.User) (error, int) {
	_, err := GetDB().Exec(context.Background(), "UPDATE users SET name = $1, surname = $2, patronymic = $3, gender = $4, date_born = $5 WHERE id_account = $6", user.Name, user.Surname, user.Patronymic, user.Gender, user.DateBorn, user.IdAccount)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err), "UserIdAccount", user.IdAccount)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
