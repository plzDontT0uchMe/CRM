package handlers

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

func Authorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "method not allowed",
		})
		w.Write(resp)
		logger.CreateLog("error", "Authorization handler, method not allowed")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "read body error",
		})
		w.Write(resp)
		logger.CreateLog("error", "Authorization handler, read body error")
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "unmarshal error",
		})
		w.Write(resp)
		logger.CreateLog("error", "Authorization handler, unmarshal error")
		return
	}

	if user.Login == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "empty login or password",
		})
		w.Write(resp)
		logger.CreateLog("error", "Authorization handler, empty login or password")
		return
	}

	conn, err := grpc.Dial(config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		w.Write(resp)
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	resp, _ := auth.Authorization(context.Background(), &authService.AuthorizationRequest{Login: user.Login, Password: user.Password})

	if resp.Successfully {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    resp.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  resp.DateExpirationAccessToken.AsTime(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			Path:     "/api/updateToken",
			HttpOnly: true,
			Expires:  resp.DateExpirationRefreshToken.AsTime(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
	}

	w.WriteHeader(int(resp.HttpStatus))
	res, _ := json.Marshal(map[string]any{
		"successfully": resp.Successfully,
		"message":      resp.Message,
	})
	w.Write(res)

	logger.CreateLog("info", resp.Message)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "method not allowed",
		})
		w.Write(resp)
		logger.CreateLog("error", "Registration handler, method not allowed")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "read body error",
		})
		w.Write(resp)
		logger.CreateLog("error", "Registration handler, read body error")
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "unmarshal error",
		})
		w.Write(resp)
		logger.CreateLog("error", "Registration handler, unmarshal error")
		return
	}

	if user.Login == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "empty login or password",
		})
		w.Write(resp)
		logger.CreateLog("error", "Registration handler, empty login or password")
		return
	}

	conn, err := grpc.Dial(config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		w.Write(resp)
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	resp, _ := auth.Registration(context.Background(), &authService.RegistrationRequest{Login: user.Login, Password: user.Password})

	if resp.Successfully {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    resp.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  resp.DateExpirationAccessToken.AsTime(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			Path:     "/api/updateToken",
			HttpOnly: true,
			Expires:  resp.DateExpirationRefreshToken.AsTime(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
	}

	w.WriteHeader(int(resp.HttpStatus))
	res, _ := json.Marshal(map[string]any{
		"successfully": resp.Successfully,
		"message":      resp.Message,
	})
	w.Write(res)

	logger.CreateLog("info", resp.Message)
}

func CheckAuthorization(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "access token not found in cookie, authorization failed",
		})
		w.Write(resp)
		logger.CreateLog("error", "CheckAuthorization handler, access token not found in cookie")
		return
	}

	var session models.Session
	session.AccessToken = cookie.Value

	conn, err := grpc.Dial(config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		w.Write(resp)
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	resp, _ := auth.CheckAuthorization(context.Background(), &authService.CheckAuthorizationRequest{AccessToken: session.AccessToken})

	w.WriteHeader(int(resp.HttpStatus))
	res, _ := json.Marshal(map[string]any{
		"successfully": resp.Successfully,
		"message":      resp.Message,
	})
	w.Write(res)

	logger.CreateLog("info", resp.Message)
}

func UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "refresh token not found in cookie, authorization failed",
		})
		w.Write(resp)
		logger.CreateLog("error", "CheckAuthorization handler, refresh token not found in cookie")
		return
	}

	var session models.Session
	session.RefreshToken = cookie.Value

	conn, err := grpc.Dial(config.GetConfig().AuthService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for auth service",
		})
		w.Write(resp)
		logger.CreateLog("error", "connection error for auth service")
		return
	}
	defer conn.Close()

	auth := authService.NewAuthServiceClient(conn)

	resp, _ := auth.UpdateAccessToken(context.Background(), &authService.UpdateAccessTokenRequest{RefreshToken: session.RefreshToken})

	if resp.Successfully {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    resp.NewAccessToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  resp.NewDateExpirationAccessToken.AsTime(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
	}

	w.WriteHeader(int(resp.HttpStatus))
	res, _ := json.Marshal(map[string]any{
		"successfully": resp.Successfully,
		"message":      resp.Message,
	})
	w.Write(res)

	logger.CreateLog("info", resp.Message)
}

func GetHelloWorld() http.Handler {
	getHelloWorld := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      "Hello world!",
		})
		w.Write(resp)
	}
	return http.HandlerFunc(getHelloWorld)
}
