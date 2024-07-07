package schedule

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"CRM/go/apiGateway/internal/proto/scheduleService"
	"CRM/go/apiGateway/internal/proto/storageService"
	"CRM/go/apiGateway/internal/proto/usersService"
	"CRM/go/apiGateway/pkg/utils"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
)

func GetRecords() http.Handler {
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
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().ScheduleService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for schedule service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for schedule service")
			return
		}
		defer conn.Close()

		schedule := scheduleService.NewScheduleServiceClient(conn)

		respSchedule, _ := schedule.GetRecords(context.Background(), &scheduleService.GetRecordsRequest{})

		if respSchedule == nil || respSchedule.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from schedule service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetRecords handler, response from schedule service is nil")
			return
		}

		if !respSchedule.Status.Successfully {
			w.WriteHeader(int(respSchedule.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSchedule.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSchedule.Status.Message)
			return
		}

		var ids []int64
		for _, recordTemp := range respSchedule.Records {
			if !utils.Contains(ids, recordTemp.ClientId) && recordTemp.ClientId != 0 {
				ids = append(ids, recordTemp.ClientId)
			}
			if !utils.Contains(ids, recordTemp.TrainerId) && recordTemp.TrainerId != 0 {
				ids = append(ids, recordTemp.TrainerId)
			}
		}

		respAuthForUsers, err := auth.GetAccounts(context.Background(), &authService.GetAccountsRequest{
			Id: ids,
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
			Id: ids,
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

		if respUsersForUsers.User[respAuth.Account.Id].Position != "manager" {
			w.WriteHeader(http.StatusForbidden)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "user is not a manager",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", "user is not a manager")
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
			Id: ids,
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

		var records []models.Record
		for _, recordTemp := range respSchedule.Records {
			client := &models.User{
				ID:           int(recordTemp.ClientId),
				LastActivity: respAuthForUsers.Accounts[recordTemp.ClientId].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[recordTemp.ClientId].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[recordTemp.ClientId].Name,
				Surname:      respUsersForUsers.User[recordTemp.ClientId].Surname,
				Patronymic:   respUsersForUsers.User[recordTemp.ClientId].Patronymic,
				Gender:       int(respUsersForUsers.User[recordTemp.ClientId].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[recordTemp.ClientId].DateBorn),
				Image:        respStorageForUsers.Links[recordTemp.ClientId],
				Position:     respUsersForUsers.User[recordTemp.ClientId].Position,
				TrainerInfo:  nil,
				Subscription: nil,
				Application:  nil,
			}

			var trainer *models.User

			if recordTemp.TrainerId != 0 {
				trainer = &models.User{
					ID:           int(recordTemp.TrainerId),
					LastActivity: respAuthForUsers.Accounts[recordTemp.TrainerId].LastActivity.AsTime(),
					DateCreated:  respAuthForUsers.Accounts[recordTemp.TrainerId].DateCreated.AsTime(),
					Name:         respUsersForUsers.User[recordTemp.TrainerId].Name,
					Surname:      respUsersForUsers.User[recordTemp.TrainerId].Surname,
					Patronymic:   respUsersForUsers.User[recordTemp.TrainerId].Patronymic,
					Gender:       int(respUsersForUsers.User[recordTemp.TrainerId].Gender),
					DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[recordTemp.TrainerId].DateBorn),
					Image:        respStorageForUsers.Links[recordTemp.TrainerId],
					Position:     respUsersForUsers.User[recordTemp.TrainerId].Position,
					TrainerInfo:  nil,
					Subscription: nil,
					Application:  nil,
				}
			} else {
				trainer = nil
			}

			record := &models.Record{
				ID:        int(recordTemp.Id),
				Client:    client,
				Trainer:   trainer,
				DateStart: recordTemp.DateStart.AsTime(),
				DateEnd:   recordTemp.DateEnd.AsTime(),
			}

			records = append(records, *record)
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      "records received successfully",
			"records":      records,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}

		logger.CreateLog("info", "records received successfully")
	})
}

func GetRecordsForUser() http.Handler {
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
			}
			logger.CreateLog("info", respAuth.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().ScheduleService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for schedule service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for schedule service")
			return
		}
		defer conn.Close()

		schedule := scheduleService.NewScheduleServiceClient(conn)

		respSchedule, _ := schedule.GetRecordsForUser(context.Background(), &scheduleService.GetRecordsForUserRequest{
			UserId: respAuth.Account.Id,
		})

		if respSchedule == nil || respSchedule.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from schedule service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetRecords handler, response from schedule service is nil")
			return
		}

		if !respSchedule.Status.Successfully {
			w.WriteHeader(int(respSchedule.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSchedule.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSchedule.Status.Message)
			return
		}

		var ids []int64
		for _, recordTemp := range respSchedule.Records {
			if !utils.Contains(ids, recordTemp.ClientId) && recordTemp.ClientId != 0 {
				ids = append(ids, recordTemp.ClientId)
			}
			if !utils.Contains(ids, recordTemp.TrainerId) && recordTemp.TrainerId != 0 {
				ids = append(ids, recordTemp.TrainerId)
			}
		}

		respAuthForUsers, err := auth.GetAccounts(context.Background(), &authService.GetAccountsRequest{
			Id: ids,
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
			Id: ids,
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
			Id: ids,
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

		var records []models.Record
		for _, recordTemp := range respSchedule.Records {
			client := &models.User{
				ID:           int(recordTemp.ClientId),
				LastActivity: respAuthForUsers.Accounts[recordTemp.ClientId].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[recordTemp.ClientId].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[recordTemp.ClientId].Name,
				Surname:      respUsersForUsers.User[recordTemp.ClientId].Surname,
				Patronymic:   respUsersForUsers.User[recordTemp.ClientId].Patronymic,
				Gender:       int(respUsersForUsers.User[recordTemp.ClientId].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[recordTemp.ClientId].DateBorn),
				Image:        respStorageForUsers.Links[recordTemp.ClientId],
				Position:     respUsersForUsers.User[recordTemp.ClientId].Position,
				TrainerInfo:  nil,
				Subscription: nil,
				Application:  nil,
			}

			var trainer *models.User

			if recordTemp.TrainerId != 0 {
				trainer = &models.User{
					ID:           int(recordTemp.TrainerId),
					LastActivity: respAuthForUsers.Accounts[recordTemp.TrainerId].LastActivity.AsTime(),
					DateCreated:  respAuthForUsers.Accounts[recordTemp.TrainerId].DateCreated.AsTime(),
					Name:         respUsersForUsers.User[recordTemp.TrainerId].Name,
					Surname:      respUsersForUsers.User[recordTemp.TrainerId].Surname,
					Patronymic:   respUsersForUsers.User[recordTemp.TrainerId].Patronymic,
					Gender:       int(respUsersForUsers.User[recordTemp.TrainerId].Gender),
					DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[recordTemp.TrainerId].DateBorn),
					Image:        respStorageForUsers.Links[recordTemp.TrainerId],
					Position:     respUsersForUsers.User[recordTemp.TrainerId].Position,
					TrainerInfo:  nil,
					Subscription: nil,
					Application:  nil,
				}
			} else {
				trainer = nil
			}

			record := &models.Record{
				ID:        int(recordTemp.Id),
				Client:    client,
				Trainer:   trainer,
				DateStart: recordTemp.DateStart.AsTime(),
				DateEnd:   recordTemp.DateEnd.AsTime(),
			}

			records = append(records, *record)
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      "records received successfully",
			"records":      records,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}

		logger.CreateLog("info", "records received successfully")
	})
}

func AddRecord() http.Handler {
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

		recordModel := models.Record{
			Client: &models.User{
				ID: int(respAuth.Account.Id),
			},
		}

		err = json.Unmarshal(b, &recordModel)
		if err != nil || &recordModel == nil {
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
			logger.CreateLog("error", "application handler, unmarshal error", err)
			return
		}

		if &recordModel.DateStart == nil || &recordModel.DateEnd == nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "invalid data",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, invalid data")
			return
		}

		record := scheduleService.Record{
			ClientId:  int64(recordModel.Client.ID),
			DateStart: timestamppb.New(recordModel.DateStart),
			DateEnd:   timestamppb.New(recordModel.DateEnd),
		}

		if recordModel.Trainer != nil {
			record.TrainerId = int64(recordModel.Trainer.ID)
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().ScheduleService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for schedule service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for schedule service")
			return
		}
		defer conn.Close()

		schedule := scheduleService.NewScheduleServiceClient(conn)

		respSchedule, _ := schedule.AddRecord(context.Background(), &scheduleService.AddRecordRequest{
			Record: &record,
		})

		if respSchedule == nil || respSchedule.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from schedule service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "AddRecord handler, response from schedule service is nil")
			return
		}

		if !respSchedule.Status.Successfully {
			w.WriteHeader(int(respSchedule.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSchedule.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSchedule.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respSchedule.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", respSchedule.Status.Message)
	})
}

func DeleteRecord() http.Handler {
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

		record := models.Record{}

		err = json.Unmarshal(b, &record)
		if err != nil || &record == nil {
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
			logger.CreateLog("error", "application handler, unmarshal error", err)
			return
		}

		if record.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "invalid data",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, invalid data")
			return
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().ScheduleService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for schedule service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for schedule service")
			return
		}
		defer conn.Close()

		schedule := scheduleService.NewScheduleServiceClient(conn)

		respSchedule, _ := schedule.DeleteRecordById(context.Background(), &scheduleService.DeleteRecordByIdRequest{
			Id: int64(record.ID),
		})

		if respSchedule == nil || respSchedule.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from schedule service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "DeleteRecord handler, response from schedule service is nil")
			return
		}

		if !respSchedule.Status.Successfully {
			w.WriteHeader(int(respSchedule.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSchedule.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSchedule.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respSchedule.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}

		logger.CreateLog("info", respSchedule.Status.Message)
	})
}

func GetRecordsByTrainerForDay() http.Handler {
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

		recordForDay := models.RecordForDay{}

		err = json.Unmarshal(b, &recordForDay)
		if err != nil || &recordForDay == nil {
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
			logger.CreateLog("error", "application handler, unmarshal error", err)
			return
		}

		if recordForDay.TrainerId == 0 || recordForDay.Day.String() == "0001-01-01 00:00:00 +0000 UTC" {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "invalid data",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "application handler, invalid data")
			return
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().ScheduleService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for schedule service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for schedule service")
			return
		}
		defer conn.Close()

		schedule := scheduleService.NewScheduleServiceClient(conn)

		respSchedule, _ := schedule.GetRecordsByTrainerForDay(context.Background(), &scheduleService.GetRecordsByTrainerForDayRequest{
			TrainerId: recordForDay.TrainerId,
			Day:       timestamppb.New(recordForDay.Day),
		})

		if respSchedule == nil || respSchedule.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from schedule service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetRecordsByTrainerForDay handler, response from schedule service is nil")
			return
		}

		if !respSchedule.Status.Successfully {
			w.WriteHeader(int(respSchedule.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSchedule.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSchedule.Status.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      "records received successfully",
			"times":        respSchedule.Time,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}

		logger.CreateLog("info", "records received successfully")
	})
}
