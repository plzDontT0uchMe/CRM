package postgres

import (
	"CRM/go/authService/internal/logger"
	"CRM/go/authService/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func GetUser(user *models.User) (error, int) {
	err := GetDB().QueryRow("SELECT id FROM users WHERE login = $1 AND password = $2 LIMIT 1", user.Login, user.Password).Scan(&user.Id)

	if errors.Is(err, sql.ErrNoRows) {
		logger.CreateLog("error", fmt.Sprintf("user not found: %v", err), "userLogin", user.Login)
		return err, http.StatusBadRequest
	}

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("query error for user: %v", err), "userLogin", user.Login)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func CreateUser(user *models.User) (error, int) {
	err := GetDB().QueryRow("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", user.Login, user.Password).Scan(&user.Id)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err), "userLogin", user.Login)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func RemoveAllSessionsByUser(user *models.User) (error, int) {
	_, err := GetDB().Exec("DELETE FROM sessions WHERE id_user = $1", user.Id)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err), "userId", user.Id, "userLogin", user.Login)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func CreateSession(session *models.Session) (error, int) {
	_, err := GetDB().Exec("INSERT INTO sessions (id_user, access_token, date_expiration_access_token, refresh_token, date_expiration_refresh_token) VALUES ($1, $2, $3, $4, $5)",
		session.UserId, session.AccessToken, session.DateExpirationAccessToken, session.RefreshToken, session.DateExpirationRefreshToken)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err), "userId", session.UserId)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func GetSessionByAccessToken(session *models.Session) (error, int) {
	err := GetDB().QueryRow("SELECT id, date_expiration_access_token FROM sessions WHERE access_token = $1 LIMIT 1", session.AccessToken).Scan(&session.Id, &session.DateExpirationAccessToken)

	if errors.Is(err, sql.ErrNoRows) {
		logger.CreateLog("error", fmt.Sprintf("session not found: %v", err))
		return err, http.StatusUnauthorized
	}

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error binding JSON for session: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func RemoveAccessTokenBySessionId(session *models.Session) (error, int) {
	_, err := GetDB().Exec("UPDATE sessions SET access_token = NULL, date_expiration_access_token = NULL WHERE id = $1", session.Id)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err))
		return err, http.StatusInternalServerError
	}

	return errors.New("access token expired"), http.StatusUnauthorized
}

func GetSessionByRefreshToken(session *models.Session) (error, int) {
	err := GetDB().QueryRow("SELECT id, date_expiration_refresh_token FROM sessions WHERE refresh_token = $1 LIMIT 1", session.RefreshToken).Scan(&session.Id, &session.DateExpirationRefreshToken)

	if errors.Is(err, sql.ErrNoRows) {
		logger.CreateLog("error", fmt.Sprintf("session not found: %v", err))
		return err, http.StatusUnauthorized
	}

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error binding JSON for session: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func RemoveSessionBySessionId(session *models.Session) (error, int) {
	_, err := db.Exec("DELETE FROM sessions WHERE id = $1", session.Id)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err))
		return err, http.StatusInternalServerError
	}

	return errors.New("refresh token expired"), http.StatusUnauthorized
}

func UpdateAccessTokenByRefreshToken(session *models.Session) (error, int) {
	_, err := GetDB().Exec("UPDATE sessions SET access_token = $1, date_expiration_access_token = $2 WHERE refresh_token = $3",
		session.AccessToken, session.DateExpirationAccessToken, session.RefreshToken)

	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("database query error: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
