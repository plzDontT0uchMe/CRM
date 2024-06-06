package handlers

import (
	"CRM/go/apiGateway/internal/config"
	"CRM/go/apiGateway/internal/logger"
	"CRM/go/apiGateway/internal/models"
	"CRM/go/apiGateway/internal/proto/authService"
	"CRM/go/apiGateway/internal/proto/storageService"
	"CRM/go/apiGateway/internal/proto/usersService"
	"CRM/go/apiGateway/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	b, err := io.ReadAll(r.Body)
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

	var account models.Account
	err = json.Unmarshal(b, &account)
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

	if account.Login == "" || account.Password == "" {
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

	respAuth, _ := auth.Authorization(context.Background(), &authService.AuthorizationRequest{
		Login:    account.Login,
		Password: account.Password,
	})

	if respAuth == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from auth service is nil",
		})
		w.Write(res)
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

	if !respAuth.Successfully {
		w.WriteHeader(int(respAuth.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": respAuth.Successfully,
			"message":      respAuth.Message,
		})
		w.Write(res)
		logger.CreateLog("info", respAuth.Message)
		return
	}

	conn, err = grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for users service",
		})
		w.Write(res)
		logger.CreateLog("error", "connection error for users service")
		return
	}
	defer conn.Close()

	users := usersService.NewUsersServiceClient(conn)

	respUsers, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{IdAccount: respAuth.IdAccount})

	if respUsers == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from users service is nil",
		})
		w.Write(res)
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

	if !respUsers.Successfully {
		w.WriteHeader(int(respUsers.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": respUsers.Successfully,
			"message":      respUsers.Message,
		})
		w.Write(res)
		logger.CreateLog("info", respUsers.Message)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    respAuth.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  respAuth.DateExpirationAccessToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    respAuth.RefreshToken,
		Path:     "/api/updateToken",
		HttpOnly: true,
		Expires:  respAuth.DateExpirationRefreshToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	user := models.User{
		ID:           int(respAuth.IdAccount),
		Role:         int(respAuth.Role),
		LastActivity: respAuth.LastActivity.AsTime(),
		DateCreated:  respAuth.DateCreated.AsTime(),
		Name:         respUsers.Name,
		Surname:      respUsers.Surname,
		Patronymic:   respUsers.Patronymic,
		Gender:       int(respUsers.Gender),
		DateBorn:     utils.ConvertTimestampToNullString(respUsers.DateBorn),
		Image:        respUsers.LinkImage,
	}
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respAuth.Message,
		"user":         user,
	})
	w.Write(res)

	logger.CreateLog("info", respAuth.Message)
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

	var account models.Account
	err = json.Unmarshal(b, &account)
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

	if account.Login == "" || account.Password == "" {
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

	respAuth, _ := auth.Registration(context.Background(), &authService.RegistrationRequest{
		Login:    account.Login,
		Password: account.Password,
	})

	if respAuth == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from auth service is nil",
		})
		w.Write(res)
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

	if !respAuth.Successfully {
		w.WriteHeader(int(respAuth.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": respAuth.Successfully,
			"message":      respAuth.Message,
		})
		w.Write(res)
		logger.CreateLog("info", respAuth.Message)
		return
	}

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
		w.Write(res)
		logger.CreateLog("error", "Registration handler, response from users service is nil")
		return
	}

	if !respUsers.Successfully {
		w.WriteHeader(int(respUsers.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": respUsers.Successfully,
			"message":      respUsers.Message,
		})
		w.Write(res)
		logger.CreateLog("info", respUsers.Message)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    respAuth.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  respAuth.DateExpirationAccessToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    respAuth.RefreshToken,
		Path:     "/api/updateToken",
		HttpOnly: true,
		Expires:  respAuth.DateExpirationRefreshToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	user := models.User{
		ID:           int(respAuth.IdAccount),
		Role:         int(respAuth.Role),
		LastActivity: respAuth.LastActivity.AsTime(),
		DateCreated:  respAuth.DateCreated.AsTime(),
		Name:         "",
		Surname:      "",
		Patronymic:   "",
		Gender:       0,
		DateBorn:     "",
		Image:        "",
	}
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      respAuth.Message,
		"user":         user,
	})
	w.Write(res)

	logger.CreateLog("info", respAuth.Message)
}

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAccount, err := strconv.Atoi(r.Header.Get("idAccount"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "idAccount not found in header",
			})
			w.Write(res)
			logger.CreateLog("error", "GetUser handler, idAccount not found in header")
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

		resp, _ := auth.Logout(context.Background(), &authService.LogoutRequest{
			IdAccount: int64(idAccount),
		})

		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from auth service is nil",
			})
			w.Write(res)
			logger.CreateLog("error", "Logout handler, response from auth service is nil")
			return
		}

		if !resp.Successfully {
			w.WriteHeader(int(resp.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": resp.Successfully,
				"message":      resp.Message,
			})
			w.Write(res)
			logger.CreateLog("info", resp.Message)
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
			"message":      resp.Message,
		})
		w.Write(res)

		logger.CreateLog("info", resp.Message)
	})
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

	resp, _ := auth.CheckAuthorization(context.Background(), &authService.CheckAuthorizationRequest{
		AccessToken: session.AccessToken,
	})

	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from auth service is nil",
		})
		w.Write(res)
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

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

	resp, _ := auth.UpdateAccessToken(context.Background(), &authService.UpdateAccessTokenRequest{
		RefreshToken: session.RefreshToken,
	})

	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from auth service is nil",
		})
		w.Write(res)
		logger.CreateLog("error", "Authorization handler, response from auth service is nil")
		return
	}

	if !resp.Successfully {
		w.WriteHeader(int(resp.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": resp.Successfully,
			"message":      resp.Message,
		})
		w.Write(res)
		logger.CreateLog("info", resp.Message)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.NewAccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  resp.NewDateExpirationAccessToken.AsTime(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]any{
		"successfully": true,
		"message":      resp.Message,
	})
	w.Write(res)

	logger.CreateLog("info", resp.Message)
}

func GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAccount, err := strconv.Atoi(r.Header.Get("idAccount"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "idAccount not found in header",
			})
			w.Write(res)
			logger.CreateLog("error", "GetUser handler, idAccount not found in header")
			return
		}
		roleAccount, err := strconv.Atoi(r.Header.Get("roleAccount"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "roleAccount not found in header",
			})
			w.Write(res)
			logger.CreateLog("error", "GetUser handler, roleAccount not found in header")
			return
		}
		lastActivityAccount, err := utils.ConvertStringToTime(r.Header.Get("lastActivityAccount"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "lastActivityAccount not found in header",
			})
			w.Write(res)
			logger.CreateLog("error", "GetUser handler, lastActivityAccount not found in header")
			return
		}
		dateCreatedAccount, err := utils.ConvertStringToTime(r.Header.Get("dateCreatedAccount"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "dateCreatedAccount not found in header",
			})
			w.Write(res)
			logger.CreateLog("error", "GetUser handler, dateCreatedAccount not found in header")
			return
		}

		conn, err := grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			w.Write(res)
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		users := usersService.NewUsersServiceClient(conn)

		resp, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{
			IdAccount: int64(idAccount),
		})

		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			w.Write(res)
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !resp.Successfully {
			w.WriteHeader(int(resp.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": resp.Successfully,
				"message":      resp.Message,
			})
			w.Write(res)
			logger.CreateLog("info", resp.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		user := models.User{
			ID:           idAccount,
			Role:         roleAccount,
			LastActivity: lastActivityAccount,
			DateCreated:  dateCreatedAccount,
			Name:         resp.Name,
			Surname:      resp.Surname,
			Patronymic:   resp.Patronymic,
			Gender:       int(resp.Gender),
			DateBorn:     utils.ConvertTimestampToNullString(resp.DateBorn),
			Image:        resp.LinkImage,
		}

		res, _ := json.Marshal(map[string]any{
			"successfully": resp.Successfully,
			"message":      resp.Message,
			"user":         user,
		})
		w.Write(res)

		logger.CreateLog("info", resp.Message)
	})
}

func GetUserById() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		parts := strings.Split(path, "/")

		if len(parts) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "Invalid path",
			})
			w.Write(res)
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
			w.Write(res)
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

		respAuth, _ := auth.GetUser(context.Background(), &authService.GetUserRequest{
			IdAccount: int64(id),
		})

		if respAuth == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from auth service is nil",
			})
			w.Write(res)
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respAuth.Successfully {
			w.WriteHeader(int(respAuth.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": respAuth.Successfully,
				"message":      respAuth.Message,
			})
			w.Write(res)
			logger.CreateLog("info", respAuth.Message)
			return
		}

		conn, err = grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			w.Write(resp)
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		users := usersService.NewUsersServiceClient(conn)

		respUsers, _ := users.GetUser(context.Background(), &usersService.GetUserRequest{
			IdAccount: int64(id),
		})
		if respUsers == nil {
			w.WriteHeader(http.StatusInternalServerError)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "response from users service is nil",
			})
			w.Write(res)
			logger.CreateLog("error", "Authorization handler, response from auth service is nil")
			return
		}

		if !respUsers.Successfully {
			w.WriteHeader(int(respUsers.HttpStatus))
			res, _ := json.Marshal(map[string]any{
				"successfully": respUsers.Successfully,
				"message":      respUsers.Message,
			})
			w.Write(res)
			logger.CreateLog("info", respUsers.Message)
			return
		}

		w.WriteHeader(http.StatusOK)
		user := models.User{
			ID:           id,
			Role:         int(respAuth.Role),
			LastActivity: respAuth.LastActivity.AsTime(),
			DateCreated:  respAuth.DateCreated.AsTime(),
			Name:         respUsers.Name,
			Surname:      respUsers.Surname,
			Patronymic:   respUsers.Patronymic,
			Gender:       int(respUsers.Gender),
			DateBorn:     utils.ConvertTimestampToNullString(respUsers.DateBorn),
			Image:        respUsers.LinkImage,
		}

		res, _ := json.Marshal(map[string]any{
			"successfully": true,
			"message":      respUsers.Message,
			"user":         user,
		})
		w.Write(res)
	})
}

func UpdateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idAccount, err := strconv.Atoi(r.Header.Get("idAccount"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "idAccount not found in header",
			})
			w.Write(res)
			logger.CreateLog("error", "GetUser handler, idAccount not found in header")
			return
		}

		gender, _ := strconv.Atoi(r.FormValue("gender")) //Проверить на валидность
		if gender < 0 || gender > 2 {
			w.WriteHeader(http.StatusBadRequest)
			res, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "invalid gender",
			})
			w.Write(res)
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
			w.Write(res)
			logger.CreateLog("error", "invalid date born")
			return
		}

		linkImage := r.FormValue("image")
		file, fileHeader, _ := r.FormFile("image")
		if file != nil {
			file.Close()

			fileFormat := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:]) //Проверить на размер изображения 5mb
			if fileFormat != "jpg" && fileFormat != "jpeg" && fileFormat != "png" && fileFormat != "gif" {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "unsupported format",
				})
				w.Write(res)
				logger.CreateLog("error", "unsupported format")
				return
			}

			fmt.Println(fileHeader.Size)

			if fileHeader.Size > 3<<20 {
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "file size exceeds 3MB",
				})
				w.Write(res)
				logger.CreateLog("error", "file size exceeds 3MB")
				return
			}

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				logger.CreateLog("error", fmt.Sprintf("error read file: %v", err))
				w.WriteHeader(http.StatusBadRequest)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "error read file",
				})
				w.Write(res)
				return
			}

			conn, err := grpc.Dial(config.GetConfig().StorageService.Address, grpc.WithInsecure())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "connection error for users service",
				})
				w.Write(resp)
				logger.CreateLog("error", "connection error for users service")
				return
			}
			defer conn.Close()

			storage := storageService.NewStorageServiceClient(conn)

			respStorage, _ := storage.UploadImage(context.Background(), &storageService.UploadImageRequest{
				Image:  fileBytes,
				Format: fileFormat,
			})

			if respStorage == nil {
				w.WriteHeader(http.StatusInternalServerError)
				res, _ := json.Marshal(map[string]any{
					"successfully": false,
					"message":      "response from users service is nil",
				})
				w.Write(res)
				logger.CreateLog("error", "Authorization handler, response from auth service is nil")
				return
			}

			if !respStorage.Successfully {
				w.WriteHeader(int(respStorage.HttpStatus))
				res, _ := json.Marshal(map[string]any{
					"successfully": respStorage.Successfully,
					"message":      respStorage.Message,
				})
				w.Write(res)
				return
			}

			linkImage = respStorage.Link
		}

		conn, err := grpc.Dial(config.GetConfig().UsersService.Address, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]any{
				"successfully": false,
				"message":      "connection error for users service",
			})
			w.Write(resp)
			logger.CreateLog("error", "connection error for users service")
			return
		}
		defer conn.Close()

		users := usersService.NewUsersServiceClient(conn)

		respUsers, _ := users.UpdateUser(context.Background(), &usersService.UpdateUserRequest{
			IdAccount:  int64(idAccount),
			Name:       r.FormValue("name"),
			Surname:    r.FormValue("surname"),
			Patronymic: r.FormValue("patronymic"),
			Gender:     int64(gender),
			DateBorn:   dateBorn,
			LinkImage:  linkImage,
		})

		w.WriteHeader(int(respUsers.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": respUsers.Successfully,
			"message":      respUsers.Message,
		})
		w.Write(res)
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
		w.Write(res)
		return
	}

	linkImage := parts[3]

	conn, err := grpc.Dial(config.GetConfig().StorageService.Address, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "connection error for storage service",
		})
		w.Write(resp)
		logger.CreateLog("error", "connection error for storage service")
		return
	}
	defer conn.Close()

	storage := storageService.NewStorageServiceClient(conn)

	resp, _ := storage.GetImage(context.Background(), &storageService.GetImageRequest{
		LinkImage: linkImage,
	})

	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		res, _ := json.Marshal(map[string]any{
			"successfully": false,
			"message":      "response from storage service is nil",
		})
		w.Write(res)
		logger.CreateLog("error", "GetImage handler, response from storage service is nil")
		return
	}

	if !resp.Successfully {
		w.WriteHeader(int(resp.HttpStatus))
		res, _ := json.Marshal(map[string]any{
			"successfully": resp.Successfully,
			"message":      resp.Message,
		})
		w.Write(res)
		logger.CreateLog("info", resp.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(resp.Image)

	logger.CreateLog("info", resp.Message)
}
