package service

import (
	"CRM/go/usersService/internal/database/postgres"
	"CRM/go/usersService/internal/models"
	"CRM/go/usersService/internal/proto/usersService"
	"CRM/go/usersService/pkg/utils"
)

func RegisterUser(r *usersService.RegistrationRequest) (error, int) {
	user := models.User{
		IdAccount: int(r.IdAccount),
		Gender:    0,
	}

	err, httpStatus := postgres.CreateUser(&user)
	if err != nil {
		return err, httpStatus
	}

	return nil, httpStatus
}

func GetUser(r *usersService.GetUserRequest) (*models.User, error, int) {
	user := models.User{
		IdAccount: int(r.IdAccount),
	}

	err, httpStatus := postgres.GetUser(&user)
	if err != nil {
		return nil, err, httpStatus
	}

	return &user, nil, httpStatus
}

func UpdateUser(r *usersService.UpdateUserRequest) (error, int) { //refactor
	user := models.User{
		IdAccount:  int(r.IdAccount),
		Name:       utils.ConvertStringToSQLNullString(r.Name),
		Surname:    utils.ConvertStringToSQLNullString(r.Surname),
		Patronymic: utils.ConvertStringToSQLNullString(r.Patronymic),
		Gender:     int(r.Gender),
		DateBorn:   utils.ConvertTimestampToSQLNullTime(r.DateBorn),
	}
	err, httpStatus := postgres.UpdateUser(user)
	if err != nil {
		return err, httpStatus
	}

	return nil, httpStatus
}
