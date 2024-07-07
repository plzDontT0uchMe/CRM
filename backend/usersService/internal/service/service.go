package service

import (
	"CRM/go/usersService/internal/database/postgres"
	"CRM/go/usersService/internal/logger"
	"CRM/go/usersService/internal/proto/usersService"
	"CRM/go/usersService/pkg/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func Registration(request *usersService.RegistrationRequest, response *usersService.RegistrationResponse) {
	user := usersService.User{
		Id: request.Id,
	}

	_, err := postgres.Registration(&user)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("register user: %v", err))
		response.Status = &usersService.Status{
			Successfully: false,
			Message:      "error registering user",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "register user successfully")
	response.Status = &usersService.Status{
		Successfully: true,
		Message:      "registering user successfully",
		HttpStatus:   http.StatusOK,
	}
}

func GetUser(request *usersService.GetUserRequest, response *usersService.GetUserResponse) {
	response.User = make(map[int64]*usersService.User)
	for _, id := range request.Id {
		var trainerInfo sql.NullString
		var dateBorn sql.NullTime
		var name, surname, patronymic sql.NullString

		user := &usersService.User{
			Id: id,
		}

		row := postgres.GetUser(user)
		err := row.Scan(&name, &surname, &patronymic, &user.Gender, &dateBorn, &user.Position, &trainerInfo)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("scan user: %v", err))
			response.Status = &usersService.Status{
				Successfully: false,
				Message:      "error getting user",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		user.DateBorn = utils.ConvertSQLNullTimeToTimestamp(dateBorn)
		trainerInfo.String = strings.ReplaceAll(trainerInfo.String, "\n", "")
		user.Name = utils.ConvertSQLNullStringToString(name)
		user.Surname = utils.ConvertSQLNullStringToString(surname)
		user.Patronymic = utils.ConvertSQLNullStringToString(patronymic)

		if trainerInfo.Valid {
			for _, v := range strings.Split(utils.ConvertSQLNullStringToString(trainerInfo), "|") {
				TrainerInfo := &usersService.TrainerInfo{
					Exp:          utils.ConvertStringToInt64(strings.Split(v, ";")[0]),
					Sport:        strings.Split(v, ";")[1],
					Achievements: strings.Split(v, ";")[2],
				}
				user.TrainerInfo = append(user.TrainerInfo, TrainerInfo)
			}
		}

		response.User[id] = user
	}

	logger.CreateLog("info", "get user successfully")
	response.Status = &usersService.Status{
		Successfully: true,
		Message:      "getting user successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetTrainers(request *usersService.GetTrainersRequest, response *usersService.GetTrainersResponse) {
	rows, err := postgres.GetTrainers()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("get trainers: %v", err))
		response.Status = &usersService.Status{
			Successfully: false,
			Message:      "error getting trainers",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		trainer := &usersService.User{}
		var trainerInfo sql.NullString
		var dateBorn sql.NullTime
		err = rows.Scan(&trainer.Id, &trainer.Name, &trainer.Surname, &trainer.Patronymic, &trainer.Gender, &dateBorn, &trainer.Position, &trainerInfo)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan error: %v", err))
			response.Status = &usersService.Status{
				Successfully: false,
				Message:      "error getting trainers",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		trainer.DateBorn = utils.ConvertSQLNullTimeToTimestamp(dateBorn)
		trainerInfo.String = strings.ReplaceAll(trainerInfo.String, "\n", "")

		if trainerInfo.Valid {
			for _, v := range strings.Split(utils.ConvertSQLNullStringToString(trainerInfo), "|") {
				TrainerInfo := &usersService.TrainerInfo{
					Exp:          utils.ConvertStringToInt64(strings.Split(v, ";")[0]),
					Sport:        strings.Split(v, ";")[1],
					Achievements: strings.Split(v, ";")[2],
				}
				trainer.TrainerInfo = append(trainer.TrainerInfo, TrainerInfo)
			}
		}

		response.Users = append(response.Users, trainer)
	}

	logger.CreateLog("info", "get trainers successfully")
	response.Status = &usersService.Status{
		Successfully: true,
		Message:      "getting trainers successfully",
		HttpStatus:   http.StatusOK,
	}
}

func UpdateUser(request *usersService.UpdateUserRequest, response *usersService.UpdateUserResponse) { //refactor
	_, err := postgres.UpdateUser(request.User)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("update user: %v", err))
		response.Status = &usersService.Status{
			Successfully: false,
			Message:      "error updating user",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "update user successfully")
	response.Status = &usersService.Status{
		Successfully: true,
		Message:      "updating user successfully",
		HttpStatus:   http.StatusOK,
	}
}
