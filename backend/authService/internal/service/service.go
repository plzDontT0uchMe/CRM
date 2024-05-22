package service

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/database/postgres"
	"CRM/go/authService/internal/models"
	"CRM/go/authService/internal/proto/authService"
	"CRM/go/authService/pkg/hash"
	"time"
)

func AuthorizeUser(r *authService.AuthorizationRequest) (*models.User, error, int) {
	user := models.User{
		Login:    r.Login,
		Password: hash.GenerateHash(r.Password + config.GetConfig().Secret),
	}

	err, httpStatus := postgres.GetUser(&user)
	if err != nil {
		return nil, err, httpStatus
	}
	return &user, nil, httpStatus
}

func RegisterUser(r *authService.RegistrationRequest) (*models.User, error, int) {
	user := models.User{
		Login:    r.Login,
		Password: hash.GenerateHash(r.Password + config.GetConfig().Secret),
	}

	err, httpStatus := postgres.CreateUser(&user)
	if err != nil {
		return nil, err, httpStatus
	}
	return &user, nil, httpStatus
}

func DeleteAllSessionsByUser(user *models.User) (error, int) {
	err, httpStatus := postgres.RemoveAllSessionsByUser(user)
	if err != nil {
		return err, httpStatus
	}
	return nil, httpStatus

}

func CreateSession(user *models.User) (*models.Session, error, int) {
	session := models.Session{
		UserId:                     user.Id,
		AccessToken:                hash.GenerateRandomHash(),
		DateExpirationAccessToken:  time.Now().UTC().Add(time.Minute),
		RefreshToken:               hash.GenerateRandomHash(),
		DateExpirationRefreshToken: time.Now().UTC().Add(time.Minute * 5),
	}

	err, httpStatus := postgres.CreateSession(&session)
	if err != nil {
		return nil, err, httpStatus
	}

	return &session, nil, httpStatus
}

func CheckAuthorization(r *authService.CheckAuthorizationRequest) (error, int) {
	session := models.Session{
		AccessToken: r.AccessToken,
	}

	err, httpStatus := postgres.GetSessionByAccessToken(&session)
	if err != nil {
		return err, httpStatus
	}

	if session.DateExpirationAccessToken.UTC().Compare(time.Now().UTC()) < 0 {
		err, httpStatus = postgres.RemoveAccessTokenBySessionId(&session)
		if err != nil {
			return err, httpStatus
		}
	}

	return nil, httpStatus
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

	return &session, nil, httpStatus
}
