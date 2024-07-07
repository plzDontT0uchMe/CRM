package auth

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"CRM/go/apiGateway/internal/proto/storageService"
	"CRM/go/apiGateway/internal/proto/subsService"
	"CRM/go/apiGateway/internal/proto/usersService"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"io"
	"net/http"
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

	conn, err = grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())
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

	users := usersService.NewUsersServiceClient(conn)

	respUsers, _ := users.Registration(context.Background(), &usersService.RegistrationRequest{
		Id: respAuth.Account.Id,
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
		logger.CreateLog("error", "Registration handler, response from users service is nil")
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

	conn, err = grpc.Dial(config.GetConfig().SubsService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for subs service",
		})
		_, err = w.Write(res)
		if err != nil {
			logger.CreateLog("error", "error writing response")
			return
		}
		logger.CreateLog("error", "connection error for subs service")
		return
	}

	subs := subsService.NewSubsServiceClient(conn)

	respSubs, _ := subs.Registration(context.Background(), &subsService.RegistrationRequest{
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
			logger.CreateLog("error", "error writing response")
			return
		}
		logger.CreateLog("error", "Registration handler, response from subs service is nil")
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
		Id: respAuth.Account.Id,
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
		logger.CreateLog("error", "Registration handler, response from storage service is nil")
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
