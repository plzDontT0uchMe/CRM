package handlers

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"CRM/go/apiGateway/internal/proto/storageService"
	"CRM/go/apiGateway/internal/proto/subsService"
	"CRM/go/apiGateway/internal/proto/trainingService"
	"CRM/go/apiGateway/internal/proto/usersService"
	"CRM/go/apiGateway/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Authorization(w http.ResponseWriter, r *http.Request) {
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
		logger.CreateLog("error", "authorization handler, error read body")
		return
	}

	var account authService.Account
	err = json.Unmarshal(b, &account)
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
		logger.CreateLog("error", "authorization handler, unmarshal error")
		return
	}

	if account.Login == "" || account.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "empty login or password",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "authorization handler, empty login or password")
		return
	}

	conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	respAuth, _ := auth.Authorization(context.Background(), &authService.AuthorizationRequest{
		Account: &account,
	})

	if respAuth == nil || respAuth.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "authService is crashed",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}
		logger.CreateLog("error", "authorization handler, authService is crashed")
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

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    respAuth.Session.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  respAuth.Session.DateExpirationAccessToken.AsTime().Add(time.Minute),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    respAuth.Session.RefreshToken,
		Path:     "/api/updateToken",
		HttpOnly: true,
		Expires:  respAuth.Session.DateExpirationRefreshToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respAuth.Status.Message,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "writing response")
		return
	}

	logger.CreateLog("info", respAuth.Status.Message)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "method not allowed",
		})
		_, err := w.Write(resp)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "Registration handler, method not allowed")
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "error read body",
		})
		_, err = w.Write(resp)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "Registration handler, error read body")
		return
	}

	var account authService.Account
	err = json.Unmarshal(b, &account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "unmarshal error",
		})
		_, err = w.Write(resp)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "Registration handler, unmarshal error")
		return
	}

	if account.Login == "" || account.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "empty login or password",
		})
		_, err = w.Write(resp)
		if err != nil {
			logger.CreateLog("error", "writing response")
		}
		logger.CreateLog("error", "Registration handler, empty login or password")
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

	respAuth, _ := auth.Registration(context.Background(), &authService.RegistrationRequest{
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
			logger.CreateLog("error", "error writing response")
			return
		}
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

	if !respAuth.Status.Successfully {
		w.WriteHeader(int(respAuth.Status.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": respAuth.Status.Successfully,
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
	/*
		conn, err = grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())

		users := usersService.NewUsersServiceClient(conn)

		respUsers, _ := users.Registration(context.Background(), &usersService.RegistrationRequest{
			IdAccount: respAuth.IdAccount,
		})

		if respUsers == nil {
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
			logger.CreateLog("error", "Registration handler, response from users service is nil")
			return
		}

		if !respUsers.Successfully {
			w.WriteHeader(int(respUsers.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": respUsers.Successfully,
				"message":      respUsers.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("info", respUsers.Message)
			return
		}

		conn, err = grpc.Dial(config.GetConfig().StorageService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for storage service",
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("error", "connection error for storage service")
			return
		}
		defer conn.Close()

		storage := storageService.NewStorageServiceClient(conn)

		respStorage, _ := storage.Registration(context.Background(), &storageService.RegistrationRequest{
			IdAccount: respAuth.IdAccount,
		})

		if respStorage == nil {
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
			logger.CreateLog("error", "Registration handler, response from storage service is nil")
			return
		}

		if !respStorage.Successfully {
			w.WriteHeader(int(respStorage.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": respStorage.Successfully,
				"message":      respStorage.Message,
			})
			_, err = w.Write(res)
			if err != nil {
				logger.CreateLog("error", "error writing response")
				return
			}
			logger.CreateLog("info", respStorage.Message)
			return
		}

	*/

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    respAuth.Session.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  respAuth.Session.DateExpirationAccessToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    respAuth.Session.RefreshToken,
		Path:     "/api/updateToken",
		HttpOnly: true,
		Expires:  respAuth.Session.DateExpirationRefreshToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)

	user := models.User{
		ID:           int(respAuth.Account.Id),
		LastActivity: respAuth.Account.LastActivity.AsTime(),
		DateCreated:  respAuth.Account.DateCreated.AsTime(),
		Name:         "",
		Surname:      "",
		Patronymic:   "",
		Gender:       0,
		DateBorn:     "",
		Image:        "",
	}

	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respAuth.Status.Message,
		"user":         user,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}

	logger.CreateLog("info", respAuth.Status.Message)
}

func Logout() http.Handler {
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

		resp, _ := auth.Logout(context.Background(), &authService.LogoutRequest{
			Session: &session,
		})

		if resp == nil || resp.Status == nil {
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

		if !resp.Status.Successfully {
			w.WriteHeader(int(resp.Status.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": resp.Status.Successfully,
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

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().UTC(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/api/updateToken",
			HttpOnly: true,
			Expires:  time.Now().UTC(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      resp.Status.Message,
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "writing response")
			return
		}

		logger.CreateLog("info", resp.Status.Message)
	})
}

func CheckAuthorization(w http.ResponseWriter, r *http.Request) {
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
		}
		logger.CreateLog("error", "CheckAuthorization handler, access token not found in cookie")
		return
	}

	var session authService.Session
	session.AccessToken = cookie.Value

	conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
		}
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	resp, _ := auth.CheckAuthorization(context.Background(), &authService.CheckAuthorizationRequest{
		Session: &session,
	})

	if resp == nil || resp.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from auth service is nil",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

	w.WriteHeader(int(resp.Status.HttpStatus))
	res, _ := json.Marshal(map[string]any{
		"successfully": resp.Status.Successfully,
		"message":      resp.Status.Message,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}

	logger.CreateLog("info", resp.Status.Message)
}

func UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "refresh token not found in cookie, authorization failed",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
		}
		logger.CreateLog("error", "CheckAuthorization handler, refresh token not found in cookie")
		return
	}

	var session authService.Session
	session.RefreshToken = cookie.Value

	conn, err := grpc.DialContext(context.Background(), config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
		}
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	resp, _ := auth.UpdateAccessToken(context.Background(), &authService.UpdateAccessTokenRequest{
		Session: &session,
	})

	if resp == nil || resp.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from auth service is nil",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
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
			logger.CreateLog("error", "error writing response")
			return
		}
		logger.CreateLog("info", resp.Status.Message)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.Session.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  resp.Session.DateExpirationAccessToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      resp.Status.Message,
	})
	_, err = w.Write(res)
	if err != nil {
		logger.CreateLog("error", "error writing response")
		return
	}

	logger.CreateLog("info", resp.Status.Message)
}

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

		respSubs, _ := subs.GetSubscriptionByAccountId(context.Background(), &subsService.GetSubscriptionByAccountIdRequest{
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

		trainerAccount := &authService.Account{}
		if respSubs.Subscription.IdTrainer != 0 {
			trainerAccount = &authService.Account{
				Id: respSubs.Subscription.IdTrainer,
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

			respAuthForTrainer, _ := auth.GetAccountById(context.Background(), &authService.GetAccountByIdRequest{
				Account: trainerAccount,
			})

			if respAuthForTrainer == nil || respAuthForTrainer.Status == nil {
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

			if !respAuthForTrainer.Status.Successfully {
				w.WriteHeader(int(respAuthForTrainer.Status.HttpStatus))
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      respAuthForTrainer.Status.Message,
				})
				_, err = w.Write(res)
				if err != nil {
					logger.CreateLog("error", "writing response")
					return
				}
				logger.CreateLog("info", respAuthForTrainer.Status.Message)
				return
			}

			trainerAccount = respAuthForTrainer.Account
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

		accounts := []int64{respAuth.Account.Id}
		if trainerAccount.Id != 0 {
			accounts = append(accounts, trainerAccount.Id)
		}

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
				logger.CreateLog("error", "error writing response")
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

		var infoForTrainer []models.TrainersInfo
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
		}

		w.WriteHeader(http.StatusOK)
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

func GetImage(w http.ResponseWriter, r *http.Request) {
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

	link := parts[3]

	conn, err := grpc.DialContext(context.Background(), config.GetConfig().StorageService.Address, grpc.WithInsecure())
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

	respStorage, _ := storage.GetImage(context.Background(), &storageService.GetImageRequest{
		Link: link,
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
		logger.CreateLog("error", "GetImage handler, response from storage service is nil")
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(respStorage.Image)

	logger.CreateLog("info", respStorage.Status.Message)
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

		nameProgram := r.FormValue("name")
		if nameProgram == "" {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "name not found in request",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		descProgram := r.FormValue("description")
		if descProgram == "" {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "description not found in request",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
			return
		}

		ids, ok := r.MultipartForm.Value["ids"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "ids not found in request",
			})
			_, err := w.Write(res)
			if err != nil {
				logger.CreateLog("error", "writing response")
				return
			}
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

		idsInt := make([]int64, len(ids))
		for i, id := range ids {
			idsInt[i], _ = strconv.ParseInt(id, 10, 64)
			fmt.Println()
		}

		respTraining, _ := training.CreateProgram(context.Background(), &trainingService.CreateProgramRequest{
			Id:          respAuth.Account.Id,
			Name:        nameProgram,
			Description: descProgram,
			Ids:         idsInt,
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

			programs[i] = models.Program{
				ID:          int(program.Id),
				IDCreator:   int(program.IdCreator),
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

func ChangeSubscription() http.Handler {
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

		idSubscription := utils.ConvertStringToNullInt64(r.FormValue("idSubscription"))

		idTrainer := utils.ConvertStringToNullInt64(r.FormValue("idTrainer"))

		month := utils.ConvertStringToNullInt64(r.FormValue("dateExpiration"))

		dateExpiration := timestamppb.New(time.Now().AddDate(0, int(month), 0))

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

		respSubs, err := subs.ChangeSubscription(context.Background(), &subsService.ChangeSubscriptionRequest{
			IdClient:       respAuth.Account.Id,
			IdSubscription: idSubscription,
			IdTrainer:      idTrainer,
			DateExpiration: dateExpiration,
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
