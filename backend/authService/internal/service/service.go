package service

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/database/postgres"
	"CRM/go/authService/internal/database/redis"
	"CRM/go/authService/internal/logger"
	"CRM/go/authService/internal/models"
	"CRM/go/authService/internal/proto/authService"
	"CRM/go/authService/pkg/hash"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func AuthorizeAccount(r *authService.AuthorizationRequest) (*models.Account, error, int) {
	account := models.Account{
		Login:    r.Login,
		Password: hash.GenerateHash(r.Password + config.GetConfig().Secret),
	}

	err, httpStatus := postgres.GetAccountByLogin(&account)
	if err != nil {
		return nil, err, httpStatus
	}

	err, httpStatus = postgres.UpdateLastActivityByAccountId(account.Id, time.Now().UTC())
	if err != nil {
		return nil, err, httpStatus
	}

	return &account, nil, httpStatus
}

func RegisterAccount(r *authService.RegistrationRequest) (*models.Account, error, int) {
	account := models.Account{
		Login:        r.Login,
		Password:     hash.GenerateHash(r.Password + config.GetConfig().Secret),
		Role:         0,
		LastActivity: time.Now().UTC(),
		DateCreated:  time.Now().UTC(),
	}

	err, httpStatus := postgres.CreateAccount(&account)
	if err != nil {
		return nil, err, httpStatus
	}
	return &account, nil, httpStatus
}

func Logout(r *authService.LogoutRequest) (error, int) {
	accessToken, err, httpStatus := postgres.DeleteSessionByAccountId(int(r.IdAccount))
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error delete all sessions by account: %v", err))
		return err, httpStatus
	}

	err = redis.Set(context.Background(), fmt.Sprintf("exp:%v", *accessToken), "", time.Second).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error delete key from redis: %v", err))
		return err, http.StatusInternalServerError
	}

	return nil, httpStatus
}

func DeleteAllSessionsByAccount(account *models.Account) (error, int) {
	err, httpStatus := postgres.DeleteAllSessionsByAccount(account)
	if err != nil {
		return err, httpStatus
	}
	return nil, httpStatus

}

func CreateSession(account *models.Account) (*models.Session, error, int) {
	session := models.Session{
		IdAccount:                  account.Id,
		AccessToken:                hash.GenerateRandomHash(),
		DateExpirationAccessToken:  time.Now().UTC().Add(time.Minute),
		RefreshToken:               hash.GenerateRandomHash(),
		DateExpirationRefreshToken: time.Now().UTC().Add(time.Minute * 5),
	}

	err, httpStatus := postgres.CreateSession(&session)
	if err != nil {
		return nil, err, httpStatus
	}

	data, err := json.Marshal(account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error marshal account: %v", err), "userLogin", account.Login)
		return nil, err, httpStatus
	}

	err = redis.Set(context.Background(), session.AccessToken, data, 0).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err), "userLogin", account.Login)
		return nil, err, httpStatus
	}

	err = redis.Set(context.Background(), fmt.Sprintf("exp:%v", session.AccessToken), "", time.Minute).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err), "userLogin", account.Login)
		return nil, err, httpStatus
	}

	return &session, nil, httpStatus
}

func CheckAuthorization(r *authService.CheckAuthorizationRequest) (*models.Account, error, int) {
	session := models.Session{
		AccessToken: r.AccessToken,
	}

	val, err := redis.Get(context.Background(), session.AccessToken).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			logger.CreateLog("info", fmt.Sprintf("key not found in redis: %v", err))
		} else {
			logger.CreateLog("error", fmt.Sprintf("error get key from redis: %v", err))
			return nil, err, http.StatusInternalServerError
		}
	} else {
		var account models.Account
		err = json.Unmarshal([]byte(val), &account)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("error unmarshal account: %v", err))
			return nil, err, http.StatusInternalServerError
		}
		account.LastActivity = time.Now().UTC()
		data, err := json.Marshal(account)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("error marshal account: %v", err))
			return nil, err, http.StatusInternalServerError
		}
		err = redis.Set(context.Background(), session.AccessToken, data, 0).Err()
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err))
			return nil, err, http.StatusInternalServerError
		}
		return &account, nil, http.StatusOK
	}

	err, httpStatus := postgres.GetSessionByAccessToken(&session)
	if err != nil {
		return nil, err, httpStatus
	}

	if session.DateExpirationAccessToken.UTC().Compare(time.Now().UTC()) < 0 {
		err, httpStatus = postgres.RemoveAccessTokenBySessionId(&session)
		if err != nil {
			return nil, err, httpStatus
		}
	}

	err, httpStatus = postgres.UpdateLastActivityByAccountId(session.IdAccount, time.Now().UTC())
	if err != nil {
		return nil, err, httpStatus
	}

	account := models.Account{
		Id: session.IdAccount,
	}
	err, httpStatus = postgres.GetAccountById(&account)
	if err != nil {
		return nil, err, httpStatus
	}

	return &account, nil, httpStatus
}

func UpdateAccessToken(r *authService.UpdateAccessTokenRequest) (*models.Session, error, int) {
	session := models.Session{
		RefreshToken:              r.RefreshToken,
		AccessToken:               hash.GenerateRandomHash(),
		DateExpirationAccessToken: time.Now().UTC().Add(time.Minute),
	}

	err, httpStatus := postgres.GetSessionByRefreshToken(&session)
	if err != nil {
		return nil, err, httpStatus
	}

	if session.DateExpirationRefreshToken.UTC().Compare(time.Now().UTC()) < 0 {
		err, httpStatus = postgres.RemoveSessionBySessionId(&session)
		if err != nil {
			return nil, err, httpStatus
		}
	}

	err, httpStatus = postgres.UpdateAccessTokenByRefreshToken(&session)
	if err != nil {
		return nil, err, httpStatus
	}

	account := models.Account{
		Id: session.IdAccount,
	}

	err, httpStatus = postgres.GetAccountById(&account)
	if err != nil {
		return nil, err, httpStatus
	}

	data, err := json.Marshal(account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error marshal account: %v", err), "userLogin", account.Login)
	}

	err = redis.Set(context.Background(), session.AccessToken, data, 0).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err), "userLogin", account.Login)
	}

	err = redis.Set(context.Background(), fmt.Sprintf("exp:%v", session.AccessToken), "", time.Minute).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err), "userLogin", account.Login)
	}

	return &session, nil, httpStatus
}

func GetUser(r *authService.GetUserRequest) (*models.Account, error, int) {
	account := models.Account{
		Id: int(r.IdAccount),
	}

	err, httpStatus := postgres.GetAccountById(&account)
	if err != nil {
		return nil, err, httpStatus
	}

	return &account, nil, httpStatus
}
