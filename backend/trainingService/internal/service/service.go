package service

import (
	"CRM/go/trainingService/internal/database/postgres"
	"CRM/go/trainingService/internal/logger"
	"CRM/go/trainingService/internal/proto/trainingService"
	"CRM/go/trainingService/pkg/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func GetExercises(request *trainingService.GetExercisesRequest, response *trainingService.GetExercisesResponse) {
	rows, err := postgres.GetExercises()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &trainingService.Status{
			Successfully: false,
			Message:      "error getting exercises",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		exercise := &trainingService.Exercise{}
		var muscles sql.NullString

		err = rows.Scan(&exercise.Id, &exercise.Name, &exercise.Description, &exercise.Image, &muscles)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
			response.Status = &trainingService.Status{
				Successfully: false,
				Message:      "error getting exercises",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		if muscles.Valid {
			exercise.Muscles = strings.Split(utils.ConvertSQLNullStringToString(muscles), ",")
		} else {
			exercise.Muscles = nil
		}

		response.Exercises = append(response.Exercises, exercise)
	}

	logger.CreateLog("info", "exercises successfully received")
	response.Status = &trainingService.Status{
		Successfully: true,
		Message:      "exercises successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func GetExerciseById(request *trainingService.GetExerciseByIdRequest, response *trainingService.GetExerciseByIdResponse) {
	response.Exercise = &trainingService.Exercise{
		Id: request.Id,
	}

	var muscles sql.NullString

	row := postgres.GetExerciseById(request.Id)
	err := row.Scan(&response.Exercise.Id, &response.Exercise.Name, &response.Exercise.Description, &response.Exercise.Image, &muscles)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &trainingService.Status{
			Successfully: false,
			Message:      "error getting exercise",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	if muscles.Valid {
		response.Exercise.Muscles = strings.Split(utils.ConvertSQLNullStringToString(muscles), ",")
	} else {
		response.Exercise.Muscles = nil
	}

	logger.CreateLog("info", "exercise successfully received")
	response.Status = &trainingService.Status{
		Successfully: true,
		Message:      "exercise successfully received",
		HttpStatus:   http.StatusOK,
	}
}

func CreateProgram(request *trainingService.CreateProgramRequest, response *trainingService.CreateProgramResponse) {
	_, err := postgres.CreateProgram(request.Id, request.Name, request.Description, request.Ids)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &trainingService.Status{
			Successfully: false,
			Message:      "error creating program",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}

	logger.CreateLog("info", "program successfully created")
	response.Status = &trainingService.Status{
		Successfully: true,
		Message:      "program successfully created",
		HttpStatus:   http.StatusOK,
	}
}

func GetProgramsByUserId(request *trainingService.GetProgramsByUserIdRequest, response *trainingService.GetProgramsByUserIdResponse) {
	rows, err := postgres.GetProgramsByUserId(request.Id)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
		response.Status = &trainingService.Status{
			Successfully: false,
			Message:      "error getting programs",
			HttpStatus:   http.StatusInternalServerError,
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exercises string
		program := &trainingService.Program{}
		err = rows.Scan(&program.Id, &program.IdCreator, &program.Name, &program.Description, &exercises)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("database scan: %v", err))
			response.Status = &trainingService.Status{
				Successfully: false,
				Message:      "error getting programs",
				HttpStatus:   http.StatusInternalServerError,
			}
			return
		}

		exercisesIds := strings.Split(exercises, ",")

		for _, id := range exercisesIds {
			row := postgres.GetExerciseById(utils.ConvertStringToInt64(id))
			exercise := &trainingService.Exercise{}
			var muscles sql.NullString
			err = row.Scan(&exercise.Id, &exercise.Name, &exercise.Description, &exercise.Image, &muscles)
			if err != nil {
				logger.CreateLog("error", fmt.Sprintf("database error: %v", err))
				response.Status = &trainingService.Status{
					Successfully: false,
					Message:      "error getting exercise",
					HttpStatus:   http.StatusInternalServerError,
				}
				return
			}

			if muscles.Valid {
				exercise.Muscles = strings.Split(utils.ConvertSQLNullStringToString(muscles), ",")
			} else {
				exercise.Muscles = nil
			}

			program.Exercises = append(program.Exercises, exercise)
		}

		response.Programs = append(response.Programs, program)
	}

	logger.CreateLog("info", "programs successfully received")
	response.Status = &trainingService.Status{
		Successfully: true,
		Message:      "programs successfully received",
		HttpStatus:   http.StatusOK,
	}
}
