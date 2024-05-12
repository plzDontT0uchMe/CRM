package service

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/database/postgres"
	"CRM/go/authService/internal/model"
	"CRM/go/authService/pkg/hash"
	"time"
)

func AuthorizeUser(user *model.User) (error, int) {

	user.Password = hash.GenerateHash(user.Password + config.GetConfig().Secret)

	err, httpStatus := postgres.GetUser(user)
	if err != nil {
		return err, httpStatus
	}
	return nil, httpStatus

}

func RegisterUser(user *model.User) (error, int) {

	user.Password = hash.GenerateHash(user.Password + config.GetConfig().Secret)

	err, httpStatus := postgres.CreateUser(user)
	if err != nil {
		return err, httpStatus
	}
	return nil, httpStatus
}

func RemoveAllSessionsByUser(user *model.User) (error, int) {

	err, httpStatus := postgres.RemoveAllSessionsByUser(user)
	if err != nil {
		return err, httpStatus
	}
	return nil, httpStatus

}

func CreateSession(user *model.User) (*model.Session, error, int) {

	var session model.Session
	session.UserId = user.Id
	session.AccessToken = hash.GenerateRandomHash()
	session.DateExpirationAccessToken = time.Now().UTC().Add(time.Minute)
	session.RefreshToken = hash.GenerateRandomHash()
	session.DateExpirationRefreshToken = time.Now().UTC().Add(time.Minute * 5)

	err, httpStatus := postgres.CreateSession(&session)
	if err != nil {
		return nil, err, httpStatus
	}

	return &session, nil, httpStatus
}

func CheckAuthorization(session *model.Session) (*string, *string, error, int) {

	if session.AccessToken != "" {
		err, httpStatus := postgres.GetSessionByAccessToken(session)
		if err != nil {
			return nil, nil, err, httpStatus
		}

		if session.DateExpirationAccessToken.UTC().Compare(time.Now().UTC()) < 0 {
			err, httpStatus = postgres.RemoveAccessTokenBySessionId(session)
			if err != nil {
				return nil, nil, err, httpStatus
			}
			message := "access token expired or not valid"
			flag := "getRefreshToken"
			return &message, &flag, nil, httpStatus
		}

		message := "authorization successful"
		return &message, nil, nil, httpStatus
	}

	if session.RefreshToken != "" {
		err, httpStatus := postgres.GetSessionByRefreshToken(session)
		if err != nil {
			return nil, nil, err, httpStatus
		}

		if session.DateExpirationRefreshToken.UTC().Compare(time.Now().UTC()) < 0 {
			err, httpStatus = postgres.RemoveSessionBySessionId(session)
			if err != nil {
				return nil, nil, err, httpStatus
			}
			message := "refresh token expired"
			flag := "authorizationFailed"
			return &message, &flag, nil, httpStatus
		}

		err, httpStatus = UpdateAccessTokenByRefreshToken(session) //refactor
		if err != nil {
			return nil, nil, err, httpStatus
		}
		message := session.AccessToken
		flag := "newAccessToken"
		return &message, &flag, nil, httpStatus
	}

	message := "authorization successful"

	return &message, nil, nil, 0
}

func UpdateAccessTokenByRefreshToken(session *model.Session) (error, int) {

	session.AccessToken = hash.GenerateRandomHash()
	session.DateExpirationAccessToken = time.Now().UTC().Add(time.Minute)

	err, httpStatus := postgres.UpdateAccessTokenByRefreshToken(session)
	if err != nil {
		return err, httpStatus
	}
	return nil, httpStatus
}
