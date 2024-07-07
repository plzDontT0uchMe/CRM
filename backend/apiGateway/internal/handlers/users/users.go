package users

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"CRM/go/apiGateway/internal/proto/storageService"
	"CRM/go/apiGateway/internal/proto/subsService"
	"CRM/go/apiGateway/internal/proto/usersService"
	"CRM/go/apiGateway/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetUser() http.Handler {
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

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().SubsService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for subs service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for subs service")
			return
		}

		subs := subsService.NewSubsServiceClient(conn)

		respSubs, _ := subs.GetSubscriptionAndApplicationByAccountId(context.Background(), &subsService.GetSubscriptionAndApplicationByAccountIdRequest{
			Id: respAuth.Account.Id,
		})

		if respSubs == nil || respSubs.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from subs service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from subs service is nil")
			return
		}

		if !respSubs.Status.Successfully {
			w.WriteHeader(int(respSubs.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSubs.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSubs.Status.Message)
			return
		}

		ids := []int64{respAuth.Account.Id}
		if respSubs.Subscription.IdTrainer != 0 && !utils.Contains(ids, respSubs.Subscription.IdTrainer) {
			ids = append(ids, respSubs.Subscription.IdTrainer)
		}
		if respSubs.Application.IdTrainer != 0 && !utils.Contains(ids, respSubs.Application.IdTrainer) {
			ids = append(ids, respSubs.Application.IdTrainer)
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
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

		auth = authService.NewAuthServiceClient(conn)

		respAuthForUsers, _ := auth.GetAccounts(context.Background(), &authService.GetAccountsRequest{
			Id: ids,
		})

		if respAuthForUsers == nil || respAuthForUsers.Status == nil {
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
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		users := usersService.NewUsersServiceClient(conn)

		respUsersForUsers, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{
			Id: ids,
		})

		if respUsersForUsers == nil || respUsersForUsers.Status == nil {
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
			logger.CreateLog("error", "GetUser handler, response from users service is nil")
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
				logger.CreateLog("error", "error writing response")
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
				"message":      "connection error for storage service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for storage service")
			return
		}
		defer conn.Close()

		storage := storageService.NewStorageServiceClient(conn)

		respStorageForUsers, _ := storage.GetLinksByIdAccounts(context.Background(), &storageService.GetLinksByIdAccountsRequest{
			Id: ids,
		})

		if respStorageForUsers == nil || respStorageForUsers.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from storage service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from storage service is nil")
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

		trainerForApp := &models.User{}
		if respSubs.Application.IdTrainer != 0 {
			trainerAppTrainerInfo := make([]models.TrainersInfo, len(respUsersForUsers.User[respSubs.Application.IdTrainer].TrainerInfo))
			fmt.Println(trainerAppTrainerInfo)
			for i, info := range respUsersForUsers.User[respSubs.Application.IdTrainer].TrainerInfo {
				trainerAppTrainerInfo[i] = models.TrainersInfo{
					Exp:          int(info.Exp),
					Sport:        info.Sport,
					Achievements: info.Achievements,
				}
			}

			trainerForApp = &models.User{
				ID:           int(respAuthForUsers.Accounts[respSubs.Application.IdTrainer].Id),
				LastActivity: respAuthForUsers.Accounts[respSubs.Application.IdTrainer].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[respSubs.Application.IdTrainer].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[respSubs.Application.IdTrainer].Name,
				Surname:      respUsersForUsers.User[respSubs.Application.IdTrainer].Surname,
				Patronymic:   respUsersForUsers.User[respSubs.Application.IdTrainer].Patronymic,
				Gender:       int(respUsersForUsers.User[respSubs.Application.IdTrainer].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[respSubs.Application.IdTrainer].DateBorn),
				Image:        respStorageForUsers.Links[respSubs.Application.IdTrainer],
				Position:     respUsersForUsers.User[respSubs.Application.IdTrainer].Position,
				TrainerInfo:  trainerAppTrainerInfo,
			}
		} else {
			trainerForApp = nil
		}

		applicationSub := models.Subscription{
			ID:             int(respSubs.Application.Subscription.Id),
			Name:           respSubs.Application.Subscription.Name,
			Price:          respSubs.Application.Subscription.Price,
			Description:    respSubs.Application.Subscription.Description,
			Possibilities:  respSubs.Application.Subscription.Possibilities,
			Trainer:        trainerForApp,
			DateExpiration: utils.ConvertInt64ToSQLNullString(respSubs.Application.Duration),
		}

		application := &models.Application{
			ID:           int(respSubs.Application.Id),
			Client:       nil,
			Subscription: applicationSub,
		}

		if respSubs.Application.Id == 0 {
			application = nil
		}

		trainerForSub := &models.User{}
		if respSubs.Subscription.IdTrainer != 0 {
			trainerSubTrainerInfo := make([]models.TrainersInfo, len(respUsersForUsers.User[respSubs.Subscription.IdTrainer].TrainerInfo))
			for i, info := range respUsersForUsers.User[respSubs.Subscription.IdTrainer].TrainerInfo {
				trainerSubTrainerInfo[i] = models.TrainersInfo{
					Exp:          int(info.Exp),
					Sport:        info.Sport,
					Achievements: info.Achievements,
				}
			}

			trainerForSub = &models.User{
				ID:           int(respAuthForUsers.Accounts[respSubs.Subscription.IdTrainer].Id),
				LastActivity: respAuthForUsers.Accounts[respSubs.Subscription.IdTrainer].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[respSubs.Subscription.IdTrainer].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[respSubs.Subscription.IdTrainer].Name,
				Surname:      respUsersForUsers.User[respSubs.Subscription.IdTrainer].Surname,
				Patronymic:   respUsersForUsers.User[respSubs.Subscription.IdTrainer].Patronymic,
				Gender:       int(respUsersForUsers.User[respSubs.Subscription.IdTrainer].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[respSubs.Subscription.IdTrainer].DateBorn),
				Image:        respStorageForUsers.Links[respSubs.Subscription.IdTrainer],
				Position:     respUsersForUsers.User[respSubs.Subscription.IdTrainer].Position,
				TrainerInfo:  trainerSubTrainerInfo,
			}
		} else {
			trainerForSub = nil
		}

		subscription := models.Subscription{
			ID:             int(respSubs.Subscription.Id),
			Name:           respSubs.Subscription.Name,
			Price:          respSubs.Subscription.Price,
			Description:    respSubs.Subscription.Description,
			Possibilities:  respSubs.Subscription.Possibilities,
			Trainer:        trainerForSub,
			DateExpiration: utils.ConvertTimestampToNullStringFull(respSubs.Subscription.DateExpiration),
		}

		userTrainerInfo := make([]models.TrainersInfo, len(respUsersForUsers.User[respAuth.Account.Id].TrainerInfo))
		for i, info := range respUsersForUsers.User[respAuth.Account.Id].TrainerInfo {
			userTrainerInfo[i] = models.TrainersInfo{
				Exp:          int(info.Exp),
				Sport:        info.Sport,
				Achievements: info.Achievements,
			}
		}

		user := models.User{
			ID:           int(respAuthForUsers.Accounts[respAuth.Account.Id].Id),
			LastActivity: respAuthForUsers.Accounts[respAuth.Account.Id].LastActivity.AsTime(),
			DateCreated:  respAuthForUsers.Accounts[respAuth.Account.Id].DateCreated.AsTime(),
			Name:         respUsersForUsers.User[respAuth.Account.Id].Name,
			Surname:      respUsersForUsers.User[respAuth.Account.Id].Surname,
			Patronymic:   respUsersForUsers.User[respAuth.Account.Id].Patronymic,
			Gender:       int(respUsersForUsers.User[respAuth.Account.Id].Gender),
			DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[respAuth.Account.Id].DateBorn),
			Image:        respStorageForUsers.Links[respAuth.Account.Id],
			Position:     respUsersForUsers.User[respAuth.Account.Id].Position,
			TrainerInfo:  userTrainerInfo,
			Subscription: &subscription,
			Application:  application,
		}

		/*var infoForTrainer []models.TrainersInfo
		var trainer *models.User

		if trainerAccount.Id != 0 {
			infoForTrainer = make([]models.TrainersInfo, len(respUsers.User[trainerAccount.Id].TrainerInfo))
			for i, info := range respUsers.User[trainerAccount.Id].TrainerInfo {
				infoForTrainer[i] = models.TrainersInfo{
					Exp:          int(info.Exp),
					Sport:        info.Sport,
					Achievements: info.Achievements,
				}
			}

			trainer = &models.User{
				ID:           int(trainerAccount.Id),
				LastActivity: trainerAccount.LastActivity.AsTime(),
				DateCreated:  trainerAccount.DateCreated.AsTime(),
				Name:         respUsers.User[trainerAccount.Id].Name,
				Surname:      respUsers.User[trainerAccount.Id].Surname,
				Patronymic:   respUsers.User[trainerAccount.Id].Patronymic,
				Gender:       int(respUsers.User[trainerAccount.Id].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsers.User[trainerAccount.Id].DateBorn),
				Image:        respStorage.Links[trainerAccount.Id],
				Position:     respUsers.User[trainerAccount.Id].Position,
				TrainerInfo:  infoForTrainer,
			}
		} else {
			trainer = nil
		}

		subscription := &models.Subscription{
			ID:             int(respSubs.Subscription.Id),
			Name:           respSubs.Subscription.Name,
			Price:          respSubs.Subscription.Price,
			Description:    respSubs.Subscription.Description,
			Possibilities:  respSubs.Subscription.Possibilities,
			Trainer:        trainer,
			DateExpiration: utils.ConvertTimestampToNullStringFull(respSubs.Subscription.DateExpiration),
		}

		trainerInfo := make([]models.TrainersInfo, len(respUsers.User[respAuth.Account.Id].TrainerInfo))
		for i, info := range respUsers.User[respAuth.Account.Id].TrainerInfo {
			trainerInfo[i] = models.TrainersInfo{
				Exp:          int(info.Exp),
				Sport:        info.Sport,
				Achievements: info.Achievements,
			}
		}

		user := models.User{
			ID:           int(respAuth.Account.Id),
			LastActivity: respAuth.Account.LastActivity.AsTime(),
			DateCreated:  respAuth.Account.DateCreated.AsTime(),
			Name:         respUsers.User[respAuth.Account.Id].Name,
			Surname:      respUsers.User[respAuth.Account.Id].Surname,
			Patronymic:   respUsers.User[respAuth.Account.Id].Patronymic,
			Gender:       int(respUsers.User[respAuth.Account.Id].Gender),
			DateBorn:     utils.ConvertTimestampToNullString(respUsers.User[respAuth.Account.Id].DateBorn),
			Image:        respStorage.Links[respAuth.Account.Id],
			Position:     respUsers.User[respAuth.Account.Id].Position,
			TrainerInfo:  trainerInfo,
			Subscription: subscription,
		}*/

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      "user successfully received",
			"user":         user,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", "user successfully received")
	})
}

func GetUserById() http.Handler { //Доделать
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for auth service",
			})
			_, err = w.Write(resp)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for auth service")
			return
		}
		defer conn.Close()

		auth := authService.NewAuthServiceClient(conn)

		account := authService.Account{
			Id: int64(id),
		}

		respAuth, _ := auth.GetAccountById(context.Background(), &authService.GetAccountByIdRequest{
			Account: &account,
		})

		if respAuth == nil || respAuth.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from auth service is nil",
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

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(resp)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		users := usersService.NewUsersServiceClient(conn)

		accounts := []int64{int64(id)}

		respUsers, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{
			Id: accounts,
		})

		if respUsers == nil || respUsers.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respUsers.Status.Successfully {
			w.WriteHeader(int(respUsers.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respUsers.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respUsers.Status.Message)
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().StorageService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for storage service",
			})
			_, err = w.Write(resp)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for storage service")
			return
		}
		defer conn.Close()

		storage := storageService.NewStorageServiceClient(conn)

		respStorage, _ := storage.GetLinksByIdAccounts(context.Background(), &storageService.GetLinksByIdAccountsRequest{
			Id: accounts,
		})

		if respStorage == nil || respStorage.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from storage service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from storage service is nil")
			return
		}

		if !respStorage.Status.Successfully {
			w.WriteHeader(int(respStorage.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respStorage.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respStorage.Status.Message)
			return
		}

		trainerInfo := make([]models.TrainersInfo, len(respUsers.User[int64(id)].TrainerInfo))
		for i, info := range respUsers.User[int64(id)].TrainerInfo {
			trainerInfo[i] = models.TrainersInfo{
				Exp:          int(info.Exp),
				Sport:        info.Sport,
				Achievements: info.Achievements,
			}
		}

		w.WriteHeader(http.StatusOK)
		user := models.User{
			ID:           id,
			LastActivity: respAuth.Account.LastActivity.AsTime(),
			DateCreated:  respAuth.Account.DateCreated.AsTime(),
			Name:         respUsers.User[int64(id)].Name,
			Surname:      respUsers.User[int64(id)].Surname,
			Patronymic:   respUsers.User[int64(id)].Patronymic,
			Gender:       int(respUsers.User[int64(id)].Gender),
			DateBorn:     utils.ConvertTimestampToNullString(respUsers.User[int64(id)].DateBorn),
			Image:        respStorage.Links[int64(id)],
			Position:     respUsers.User[int64(id)].Position,
			TrainerInfo:  trainerInfo,
			Subscription: nil,
		}

		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respUsers.Status.Message,
			"user":         user,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respUsers.Status.Message)
	})
}

func GetUsersByTrainer() http.Handler {
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

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}

		users := usersService.NewUsersServiceClient(conn)

		respUsers, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{
			Id: []int64{respAuth.Account.Id},
		})

		if respUsers == nil || respUsers.Status == nil {
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
			logger.CreateLog("error", "GetUser handler, response from users service is nil")
			return
		}

		if !respUsers.Status.Successfully {
			w.WriteHeader(int(respUsers.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respUsers.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("info", respUsers.Status.Message)
			return
		}

		if respUsers.User[respAuth.Account.Id].Position != "trainer" {
			w.WriteHeader(http.StatusForbidden)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "user is not a trainer",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", "user is not a trainer")
			return
		}

		conn, err = grpc.DialContext(context.Background(), config.GetConfig().SubsService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for subs service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for subs service")
			return
		}

		subs := subsService.NewSubsServiceClient(conn)

		respSubs, _ := subs.GetUsersByTrainerId(context.Background(), &subsService.GetUsersByTrainerIdRequest{
			Id: respAuth.Account.Id,
		})

		if respSubs == nil || respSubs.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from subs service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "GetUser handler, response from subs service is nil")
			return
		}

		if !respSubs.Status.Successfully {
			w.WriteHeader(int(respSubs.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSubs.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("info", respSubs.Status.Message)
			return
		}

		var ids []int64
		for _, id := range respSubs.Id {
			if id != 0 && !utils.Contains(ids, id) {
				ids = append(ids, id)
			}
		}

		respAuthForUsers, _ := auth.GetAccounts(context.Background(), &authService.GetAccountsRequest{
			Id: ids,
		})

		if respAuthForUsers == nil || respAuthForUsers.Status == nil {
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
				"message":      "connection error for users service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}

		respUsersForUsers, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{
			Id: ids,
		})

		if respUsersForUsers == nil || respUsersForUsers.Status == nil {
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
			logger.CreateLog("error", "GetUser handler, response from users service is nil")
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
				"message":      "connection error for storage service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "connection error for storage service")
			return
		}

		storage := storageService.NewStorageServiceClient(conn)

		respStorageForUsers, _ := storage.GetLinksByIdAccounts(context.Background(), &storageService.GetLinksByIdAccountsRequest{
			Id: ids,
		})

		if respStorageForUsers == nil || respStorageForUsers.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from storage service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "Authorization handler, response from storage service is nil")
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

		usersResp := make([]models.User, len(respSubs.Id))
		for i, id := range respSubs.Id {
			usersResp[i] = models.User{
				ID:           int(respAuthForUsers.Accounts[id].Id),
				LastActivity: respAuthForUsers.Accounts[id].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[id].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[id].Name,
				Surname:      respUsersForUsers.User[id].Surname,
				Patronymic:   respUsersForUsers.User[id].Patronymic,
				Gender:       int(respUsersForUsers.User[id].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[id].DateBorn),
				Image:        respStorageForUsers.Links[id],
				Position:     respUsersForUsers.User[id].Position,
				TrainerInfo:  nil,
				Subscription: nil,
				Application:  nil,
			}
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      "users successfully received",
			"clients":      usersResp,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", "users successfully received")
	})
}

func UpdateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
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
			logger.CreateLog("error", "GetUser handler, idAccount not found in header")
			return
		}

		gender, _ := strconv.Atoi(r.FormValue("gender"))
		if gender < 0 || gender > 2 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "invalid gender",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "invalid gender")
			return
		}

		dateBorn, _ := utils.ConvertStringToTimestamp(r.FormValue("dateBorn"))
		years := int(time.Since(dateBorn.AsTime()).Hours() / 24 / 365)
		if years < 18 || years > 100 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "invalid date born",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("error", "invalid date born")
			return
		}

		file, fileHeader, _ := r.FormFile("image")
		if file != nil {
			file.Close()

			fileFormat := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
			if fileFormat != "jpg" && fileFormat != "jpeg" && fileFormat != "png" && fileFormat != "gif" {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "unsupported format",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				logger.CreateLog("error", "unsupported format")
				return
			}

			if fileHeader.Size > 2<<20 {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "file size exceeds 2MB",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				logger.CreateLog("error", "file size exceeds 2MB")
				return
			}

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "error read file",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				logger.CreateLog("error", fmt.Sprintf("read file: %v", err))
				return
			}

			conn, err := grpc.DialContext(context.Background(), config.GetConfig().StorageService.Address, grpc.WithInsecure())
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

			storage := storageService.NewStorageServiceClient(conn)

			respStorage, _ := storage.UploadImage(context.Background(), &storageService.UploadImageRequest{
				Id:     int64(id),
				Image:  fileBytes,
				Format: fileFormat,
			})

			if respStorage == nil || respStorage.Status == nil {
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

			if !respStorage.Status.Successfully {
				w.WriteHeader(int(respStorage.Status.HttpStatus))
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      respStorage.Status.Message,
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "error writing response")
					return
				}
				logger.CreateLog("info", respStorage.Status.Message)
				return
			}
		}

		conn, err := grpc.DialContext(context.Background(), config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			_, err = w.Write(resp)
			if err != nil {
				logger.CreateLog("error", "writing response")
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		users := usersService.NewUsersServiceClient(conn)

		user := usersService.User{
			Id:         int64(id),
			Name:       r.FormValue("name"),
			Surname:    r.FormValue("surname"),
			Patronymic: r.FormValue("patronymic"),
			Gender:     int64(gender),
			DateBorn:   dateBorn,
		}

		respUsers, _ := users.UpdateUser(context.Background(), &usersService.UpdateUserRequest{
			User: &user,
		})

		if respUsers == nil || respUsers.Status == nil {
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

		if !respUsers.Status.Successfully {
			w.WriteHeader(int(respUsers.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respUsers.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respUsers.Status.Message)
			return
		}

		w.WriteHeader(int(respUsers.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respUsers.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respUsers.Status.Message)
	})
}

func GetTrainers(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.DialContext(context.Background(), config.GetConfig().UsersService.Address, grpc.WithInsecure())
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

	users := usersService.NewUsersServiceClient(conn)

	respUsers, _ := users.GetTrainers(context.Background(), &usersService.GetTrainersRequest{})

	if respUsers == nil || respUsers.Status == nil {
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
		logger.CreateLog("error", "GetTrainers handler, response from users service is nil")
		return
	}

	if !respUsers.Status.Successfully {
		w.WriteHeader(int(respUsers.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      respUsers.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", respUsers.Status.Message)
		return
	}

	conn, err = grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
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

	var accounts []int64
	for _, user := range respUsers.Users {
		accounts = append(accounts, user.Id)
	}

	respAuth, _ := auth.GetAccounts(context.Background(), &authService.GetAccountsRequest{
		Id: accounts,
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

	conn, err = grpc.DialContext(context.Background(), config.GetConfig().StorageService.Address, grpc.WithInsecure())
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

	storage := storageService.NewStorageServiceClient(conn)

	respStorage, _ := storage.GetLinksByIdAccounts(context.Background(), &storageService.GetLinksByIdAccountsRequest{
		Id: accounts,
	})

	if respStorage == nil || respStorage.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from storage service is nil",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "Authorization handler, response from storage service is nil")
		return
	}

	if !respStorage.Status.Successfully {
		w.WriteHeader(int(respStorage.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      respStorage.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", respStorage.Status.Message)
		return
	}

	var trainers []models.User
	for _, user := range respUsers.Users {
		trainerInfo := make([]models.TrainersInfo, len(user.TrainerInfo))
		for i, info := range user.TrainerInfo {
			trainerInfo[i] = models.TrainersInfo{
				Exp:          int(info.Exp),
				Sport:        info.Sport,
				Achievements: info.Achievements,
			}
		}
		trainer := models.User{
			ID:           int(user.Id),
			LastActivity: respAuth.Accounts[user.Id].LastActivity.AsTime(),
			DateCreated:  respAuth.Accounts[user.Id].DateCreated.AsTime(),
			Name:         user.Name,
			Surname:      user.Surname,
			Patronymic:   user.Patronymic,
			Gender:       int(user.Gender),
			DateBorn:     utils.ConvertTimestampToNullString(user.DateBorn),
			Image:        respStorage.Links[user.Id],
			Position:     user.Position,
			TrainerInfo:  trainerInfo,
		}
		trainers = append(trainers, trainer)
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respUsers.Status.Message,
		"trainers":     trainers,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}

	logger.CreateLog("info", respUsers.Status.Message)
}
