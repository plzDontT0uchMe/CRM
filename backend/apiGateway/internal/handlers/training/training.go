package training

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"CRM/go/apiGateway/internal/proto/storageService"
	"CRM/go/apiGateway/internal/proto/trainingService"
	"CRM/go/apiGateway/internal/proto/usersService"
	"CRM/go/apiGateway/pkg/utils"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetExercises(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for exercises service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "connection error for exercises service")
		return
	}
	defer conn.Close()

	training := trainingService.NewTrainingServiceClient(conn)

	respTraining, _ := training.GetExercises(context.Background(), &trainingService.GetExercisesRequest{})

	if respTraining == nil || respTraining.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from exercises service is nil",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "GetExercises handler, response from exercises service is nil")
		return
	}

	if !respTraining.Status.Successfully {
		w.WriteHeader(int(respTraining.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", respTraining.Status.Message)
		return
	}

	exercises := make([]models.Exercise, len(respTraining.Exercises))
	for i, exercise := range respTraining.Exercises {
		muscles := &exercise.Muscles
		if len(*muscles) == 0 {
			muscles = nil
		}
		exercises[i] = models.Exercise{
			ID:          int(exercise.Id),
			Name:        exercise.Name,
			Description: exercise.Description,
			Image:       exercise.Image,
			Muscles:     muscles,
		}
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respTraining.Status.Message,
		"exercises":    exercises,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}

	logger.CreateLog("info", respTraining.Status.Message)
}

func GetExerciseById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "Invalid path",
		})
		_, err := w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		return
	}

	idStr := parts[3]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "Invalid id",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		return
	}

	conn, err := grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for exercises service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "connection error for exercises service")
		return
	}
	defer conn.Close()

	training := trainingService.NewTrainingServiceClient(conn)

	respTraining, _ := training.GetExerciseById(context.Background(), &trainingService.GetExerciseByIdRequest{
		Id: int64(id),
	})

	if respTraining == nil || respTraining.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from exercises service is nil",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "GetExercisesByMuscle handler, response from exercises service is nil")
		return
	}

	if !respTraining.Status.Successfully {
		w.WriteHeader(int(respTraining.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", respTraining.Status.Message)
		return
	}

	muscles := &respTraining.Exercise.Muscles
	if len(*muscles) == 0 {
		muscles = nil
	}
	exercise := models.Exercise{ //Добавить video
		ID:          int(respTraining.Exercise.Id),
		Name:        respTraining.Exercise.Name,
		Description: respTraining.Exercise.Description,
		Image:       respTraining.Exercise.Image,
		Muscles:     muscles,
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respTraining.Status.Message,
		"exercise":     exercise,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}

	logger.CreateLog("info", respTraining.Status.Message)
}

func CreateProgram() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "method not allowed",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "authorization handler, method not allowed")
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, access token not found in cookie")
			return
		}

		session := authService.Session{
			AccessToken: cookie.Value,
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		respAuth, _ := auth.GetAccountByAccessToken(context.Background(), &authService.GetAccountByAccessTokenRequest{
			Session: &session,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Status.Successfully {
			w.WriteHeader(int(respAuth.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuth.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "error read body",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, error read body")
			return
		}

		program := trainingService.Program{
			IdCreator: respAuth.Account.Id,
		}

		err = json.Unmarshal(b, &program)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "unmarshal error",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "CreateProgram handler, unmarshal error")
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}
		defer conn.Close()

		training := trainingService.NewTrainingServiceClient(conn)

		respTraining, _ := training.CreateProgram(context.Background(), &trainingService.CreateProgramRequest{
			Id:          respAuth.Account.Id,
			Name:        program.Name,
			Description: program.Description,
			Exercises:   program.Exercises,
		})

		if respTraining == nil || respTraining.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "CreateProgram handler, response from exercises service is nil")
			return
		}

		if !respTraining.Status.Successfully {
			w.WriteHeader(int(respTraining.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respTraining.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respTraining.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respTraining.Status.Message)
	})
}

func GetProgramsByUserId() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, access token not found in cookie")
			return
		}

		session := authService.Session{
			AccessToken: cookie.Value,
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		respAuth, _ := auth.GetAccountByAccessToken(context.Background(), &authService.GetAccountByAccessTokenRequest{
			Session: &session,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Status.Successfully {
			w.WriteHeader(int(respAuth.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuth.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		training := trainingService.NewTrainingServiceClient(conn)

		respTraining, err := training.GetProgramsByUserId(context.Background(), &trainingService.GetProgramsByUserIdRequest{
			Id: respAuth.Account.Id,
		})

		if respTraining == nil || respTraining.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetPrograms handler, response from exercises service is nil")
			return
		}

		if !respTraining.Status.Successfully {
			w.WriteHeader(int(respTraining.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respTraining.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respTraining.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		var idsCreators []int64
		for _, program := range respTraining.Programs {
			if !utils.Contains(idsCreators, program.IdCreator) {
				idsCreators = append(idsCreators, program.IdCreator)
			}
		}

		respAuthForUsers, err := auth.GetAccounts(context.Background(), &authService.GetAccountsRequest{
			Id: idsCreators,
		})

		if respAuthForUsers == nil || respAuthForUsers.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetPrograms handler, response from exercises service is nil")
			return
		}

		if !respAuthForUsers.Status.Successfully {
			w.WriteHeader(int(respAuthForUsers.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuthForUsers.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuthForUsers.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		users := usersService.NewUsersServiceClient(conn)

		respUsersForUsers, err := users.GetUser(context.Background(), &usersService.GetUserRequest{
			Id: idsCreators,
		})

		if respUsersForUsers == nil || respUsersForUsers.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetPrograms handler, response from exercises service is nil")
			return
		}

		if !respUsersForUsers.Status.Successfully {
			w.WriteHeader(int(respUsersForUsers.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respUsersForUsers.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respUsersForUsers.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().StorageService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		storage := storageService.NewStorageServiceClient(conn)

		respStorageForUsers, err := storage.GetLinksByIdAccounts(context.Background(), &storageService.GetLinksByIdAccountsRequest{
			Id: idsCreators,
		})

		if respStorageForUsers == nil || respStorageForUsers.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetPrograms handler, response from exercises service is nil")
			return
		}

		if !respStorageForUsers.Status.Successfully {
			w.WriteHeader(int(respStorageForUsers.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respStorageForUsers.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respStorageForUsers.Status.Message)
			return
		}

		programs := make([]models.Program, len(respTraining.Programs))
		for i, program := range respTraining.Programs {
			exercises := make([]models.Exercise, len(program.Exercises))
			for j, exercise := range program.Exercises {
				exercises[j] = models.Exercise{
					ID:          int(exercise.Id),
					Name:        exercise.Name,
					Description: exercise.Description,
					Image:       exercise.Image,
					Muscles:     &exercise.Muscles,
				}
			}

			user := models.User{
				ID:           int(program.IdCreator),
				LastActivity: respAuthForUsers.Accounts[program.IdCreator].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[program.IdCreator].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[program.IdCreator].Name,
				Surname:      respUsersForUsers.User[program.IdCreator].Surname,
				Patronymic:   respUsersForUsers.User[program.IdCreator].Patronymic,
				Gender:       int(respUsersForUsers.User[program.IdCreator].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[program.IdCreator].DateBorn),
				Image:        respStorageForUsers.Links[program.IdCreator],
				Position:     respUsersForUsers.User[program.IdCreator].Position,
			}

			programs[i] = models.Program{
				ID:          int(program.Id),
				Creator:     user,
				Name:        program.Name,
				Description: program.Description,
				Exercises:   exercises,
			}
		}

		w.WriteHeader(http.StatusOK)
		var res []byte
		if len(programs) == 0 {
			res, _ = json.Marshal(map[string]any{
				"successfully": true,
				"message":      respTraining.Status.Message,
			})
		} else {
			res, _ = json.Marshal(map[string]any{
				"successfully": true,
				"message":      respTraining.Status.Message,
				"programs":     programs,
			})
		}
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respTraining.Status.Message)
	})
}

func DeleteProgramLocal() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "method not allowed",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "authorization handler, method not allowed")
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, access token not found in cookie")
			return
		}

		session := authService.Session{
			AccessToken: cookie.Value,
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		respAuth, _ := auth.GetAccountByAccessToken(context.Background(), &authService.GetAccountByAccessTokenRequest{
			Session: &session,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Status.Successfully {
			w.WriteHeader(int(respAuth.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuth.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		program := trainingService.Program{}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "error read body",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, error read body")
			return
		}

		err = json.Unmarshal(b, &program)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "unmarshal error",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "CreateProgram handler, unmarshal error")
			return
		}

		if program.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "id not found in request",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		training := trainingService.NewTrainingServiceClient(conn)

		respTraining, _ := training.DeleteProgramLocal(context.Background(), &trainingService.DeleteProgramLocalRequest{
			Id:        program.Id,
			IdCreator: respAuth.Account.Id,
		})

		if respTraining == nil || respTraining.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "DeleteProgram handler, response from exercises service is nil")
			return
		}

		if !respTraining.Status.Successfully {
			w.WriteHeader(int(respTraining.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respTraining.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respTraining.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respTraining.Status.Message)
	})
}

func DeleteProgram() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "method not allowed",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "authorization handler, method not allowed")
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, access token not found in cookie")
			return
		}

		session := authService.Session{
			AccessToken: cookie.Value,
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		respAuth, _ := auth.GetAccountByAccessToken(context.Background(), &authService.GetAccountByAccessTokenRequest{
			Session: &session,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Status.Successfully {
			w.WriteHeader(int(respAuth.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuth.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		program := trainingService.Program{}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "error read body",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, error read body")
			return
		}

		err = json.Unmarshal(b, &program)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "unmarshal error",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "CreateProgram handler, unmarshal error")
			return
		}

		if program.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "id not found in request",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		training := trainingService.NewTrainingServiceClient(conn)

		respTraining, _ := training.DeleteProgram(context.Background(), &trainingService.DeleteProgramRequest{
			Id:        program.Id,
			IdCreator: respAuth.Account.Id,
		})

		if respTraining == nil || respTraining.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "DeleteProgram handler, response from exercises service is nil")
			return
		}

		if !respTraining.Status.Successfully {
			w.WriteHeader(int(respTraining.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respTraining.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respTraining.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respTraining.Status.Message)
	})
}

func ShareProgram() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "method not allowed",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "authorization handler, method not allowed")
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, access token not found in cookie")
			return
		}

		session := authService.Session{
			AccessToken: cookie.Value,
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		respAuth, _ := auth.GetAccountByAccessToken(context.Background(), &authService.GetAccountByAccessTokenRequest{
			Session: &session,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Status.Successfully {
			w.WriteHeader(int(respAuth.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuth.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		program := trainingService.Program{}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "error read body",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, error read body")
			return
		}

		err = json.Unmarshal(b, &program)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "unmarshal error",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "CreateProgram handler, unmarshal error")
			return
		}

		if program.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "id not found in request",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		if program.IdClient == 0 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "idClient not found in request",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		training := trainingService.NewTrainingServiceClient(conn)

		respTraining, _ := training.ShareProgram(context.Background(), &trainingService.ShareProgramRequest{
			Id:        program.Id,
			IdCreator: respAuth.Account.Id,
			IdClient:  program.IdClient,
		})

		if respTraining == nil || respTraining.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "DeleteProgram handler, response from exercises service is nil")
			return
		}

		if !respTraining.Status.Successfully {
			w.WriteHeader(int(respTraining.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respTraining.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respTraining.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respTraining.Status.Message)
	})
}

func ChangeProgram() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "method not allowed",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "authorization handler, method not allowed")
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "access token not found in cookie, authorization failed",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, access token not found in cookie")
			return
		}

		session := authService.Session{
			AccessToken: cookie.Value,
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		respAuth, _ := auth.GetAccountByAccessToken(context.Background(), &authService.GetAccountByAccessTokenRequest{
			Session: &session,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Status.Successfully {
			w.WriteHeader(int(respAuth.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respAuth.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "error read body",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, error read body")
			return
		}

		program := trainingService.Program{
			IdCreator: respAuth.Account.Id,
		}

		err = json.Unmarshal(b, &program)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "unmarshal error",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "CreateProgram handler, unmarshal error")
			return
		}

		if program.Id == 0 || program.Name == "" || program.Description == "" || len(program.Exercises) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "id, name, description or exercises not found in request",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().TrainingService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for exercises service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for exercises service")
			return
		}

		training := trainingService.NewTrainingServiceClient(conn)

		respTraining, _ := training.ChangeProgram(context.Background(), &trainingService.ChangeProgramRequest{
			Program: &program,
		})
		if respTraining == nil || respTraining.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from exercises service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "ChangeProgram handler, response from exercises service is nil")
			return
		}

		if !respTraining.Status.Successfully {
			w.WriteHeader(int(respTraining.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respTraining.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respTraining.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respTraining.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respTraining.Status.Message)
	})
}
