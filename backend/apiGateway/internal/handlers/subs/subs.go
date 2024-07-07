package subs

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
	"google.golang.org/grpc"
	"io"
	"net/http"
)

func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(config.GetConfig().SubsService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for storage service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "connection error for storage service")
		return
	}
	defer conn.Close()

	subs := subsService.NewSubsServiceClient(conn)

	resp, _ := subs.GetSubscriptions(context.Background(), &subsService.GetSubscriptionsRequest{})

	if resp == nil || resp.Status == nil {
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
		logger.CreateLog("error", "GetSubscriptions handler, response from storage service is nil")
		return
	}

	if !resp.Status.Successfully {
		w.WriteHeader(int(resp.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      resp.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("info", resp.Status.Message)
		return
	}

	subscriptions := make([]models.Subscription, len(resp.Subscriptions))
	for i, sub := range resp.Subscriptions {
		subscriptions[i] = models.Subscription{
			ID:            int(sub.Id),
			Name:          sub.Name,
			Price:         sub.Price,
			Description:   sub.Description,
			Possibilities: sub.Possibilities,
			Trainer:       nil,
		}
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully":  true,
		"message":       resp.Status.Message,
		"subscriptions": subscriptions,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}
}

func CreateApplication() http.Handler {
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

		application := subsService.Application{
			IdClient: respAuth.Account.Id,
		}

		err = json.Unmarshal(b, &application)
		if err != nil || application.Subscription == nil {
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
			logger.CreateLog("error", "application handler, unmarshal error")
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
			}
			logger.CreateLog("error", "connection error for subs service")
			return

		}

		subs := subsService.NewSubsServiceClient(conn)

		respSubsForCheckSub, err := subs.GetSubscriptionByAccountId(context.Background(), &subsService.GetSubscriptionByAccountIdRequest{
			Id: application.IdClient,
		})

		if respSubsForCheckSub == nil || respSubsForCheckSub.Status == nil {
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

		if !respSubsForCheckSub.Status.Successfully {
			w.WriteHeader(int(respSubsForCheckSub.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSubsForCheckSub.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSubsForCheckSub.Status.Message)
			return
		}

		if respSubsForCheckSub.Subscription.Id == application.Subscription.Id && respSubsForCheckSub.Subscription.IdTrainer == application.IdTrainer {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "subscription already exists",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		respSubsForSub, err := subs.GetSubscriptionById(context.Background(), &subsService.GetSubscriptionByIdRequest{
			Id: application.Subscription.Id,
		})

		if respSubsForSub == nil || respSubsForSub.Status == nil {
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

		if !respSubsForSub.Status.Successfully {
			w.WriteHeader(int(respSubsForSub.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      respSubsForSub.Status.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			logger.CreateLog("info", respSubsForSub.Status.Message)
			return
		}

		if respSubsForSub.Subscription.Price != 0 {
			if application.Duration != 3 && application.Duration != 6 && application.Duration != 9 && application.Duration != 12 {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "invalid duration",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				return
			}
		} else {
			application.Duration = 0
		}

		if application.IdTrainer != 0 {
			index := utils.IndexOfStrings(respSubsForSub.Subscription.Possibilities, "Sign up for a gym with a trainer")
			if index == -1 {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "subscription does not have the possibility to sign up for a gym with a trainer",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				return
			}

			conn, err = grpc.DialContext(context.Background(), config.GetConfig().UsersService.Address, grpc.WithInsecure())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "connection error for subs service",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
				}
				logger.CreateLog("error", "connection error for subs service")
				return

			}

			users := usersService.NewUsersServiceClient(conn)

			ids := []int64{application.IdTrainer}
			respUsers, err := users.GetUser(context.Background(), &usersService.GetUserRequest{
				Id: ids,
			})

			if respUsers == nil || respUsers.Status == nil {
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

			if respUsers.User[application.IdTrainer].Position != "trainer" {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "user is not a trainer",
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				return
			}
		}

		respSubs, err := subs.CreateApplication(context.Background(), &subsService.CreateApplicationRequest{
			Application: &application,
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

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respSubs.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respSubs.Status.Message)
	})
}

func GetApplications() http.Handler {
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
			}
			logger.CreateLog("error", "connection error for subs service")
			return
		}

		subs := subsService.NewSubsServiceClient(conn)

		respSubs, err := subs.GetApplications(context.Background(), &subsService.GetApplicationsRequest{})

		if respSubs == nil || respSubs.Status == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from subs service is nil",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
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
			}
			logger.CreateLog("info", respSubs.Status.Message)
			return
		}

		ids := []int64{respAuth.Account.Id}
		for _, application := range respSubs.Applications {
			if !utils.Contains(ids, application.IdClient) && application.IdClient != 0 {
				ids = append(ids, application.IdClient)
			}
			if !utils.Contains(ids, application.IdTrainer) && application.IdTrainer != 0 {
				ids = append(ids, application.IdTrainer)
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

		var applications []models.Application
		for _, application := range respSubs.Applications {
			client := models.User{
				ID:           int(application.IdClient),
				LastActivity: respAuthForUsers.Accounts[application.IdClient].LastActivity.AsTime(),
				DateCreated:  respAuthForUsers.Accounts[application.IdClient].DateCreated.AsTime(),
				Name:         respUsersForUsers.User[application.IdClient].Name,
				Surname:      respUsersForUsers.User[application.IdClient].Surname,
				Patronymic:   respUsersForUsers.User[application.IdClient].Patronymic,
				Gender:       int(respUsersForUsers.User[application.IdClient].Gender),
				DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[application.IdClient].DateBorn),
				Image:        respStorageForUsers.Links[application.IdClient],
				Position:     respUsersForUsers.User[application.IdClient].Position,
			}

			var trainer *models.User
			if application.IdTrainer != 0 {
				trainer = &models.User{
					ID:           int(application.IdTrainer),
					LastActivity: respAuthForUsers.Accounts[application.IdTrainer].LastActivity.AsTime(),
					DateCreated:  respAuthForUsers.Accounts[application.IdTrainer].DateCreated.AsTime(),
					Name:         respUsersForUsers.User[application.IdTrainer].Name,
					Surname:      respUsersForUsers.User[application.IdTrainer].Surname,
					Patronymic:   respUsersForUsers.User[application.IdTrainer].Patronymic,
					Gender:       int(respUsersForUsers.User[application.IdTrainer].Gender),
					DateBorn:     utils.ConvertTimestampToNullString(respUsersForUsers.User[application.IdTrainer].DateBorn),
					Image:        respStorageForUsers.Links[application.IdTrainer],
					Position:     respUsersForUsers.User[application.IdTrainer].Position,
				}
			} else {
				trainer = nil
			}

			subscription := models.Subscription{
				ID:             int(application.Subscription.Id),
				Name:           application.Subscription.Name,
				Price:          application.Subscription.Price,
				Description:    application.Subscription.Description,
				Possibilities:  application.Subscription.Possibilities,
				Trainer:        trainer,
				DateExpiration: utils.ConvertInt64ToSQLNullString(application.Duration),
			}

			applications = append(applications, models.Application{
				ID:           int(application.Id),
				Client:       &client,
				Subscription: subscription,
			})
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respSubs.Status.Message,
			"applications": applications,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respSubs.Status.Message)
	})
}

func ChangeSubscription() http.Handler { //добавить проверку что с акка manager
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
			}
			logger.CreateLog("error", "connection error for users service")
			return
		}

		users := usersService.NewUsersServiceClient(conn)

		respUsers, err := users.GetUser(context.Background(), &usersService.GetUserRequest{
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
			}
			logger.CreateLog("error", "Authorization handler, response from users service is nil")
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
			}
			logger.CreateLog("info", respUsers.Status.Message)
			return
		}

		if respUsers.User[respAuth.Account.Id].Position != "manager" {
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

		application := models.ChangeApplication{}

		err = json.Unmarshal(b, &application)
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
			logger.CreateLog("error", "ChangeSubscription handler, unmarshal error")
			return
		}

		if application.Application == nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "application not found",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
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
			}
			logger.CreateLog("error", "connection error for subs service")
			return
		}
		defer conn.Close()

		subs := subsService.NewSubsServiceClient(conn)

		respSubs, err := subs.ChangeApplication(context.Background(), &subsService.ChangeApplicationRequest{
			Id:         int64(application.Application.ID),
			IsAccepted: application.IsAccepted,
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

		w.WriteHeader(int(respSubs.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respSubs.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}

		logger.CreateLog("info", respSubs.Status.Message)
	})
}
