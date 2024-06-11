package service

import (
	"CRM/go/authService/internal/config"
	"CRM/go/authService/internal/database/postgres"
	"CRM/go/authService/internal/database/redis"
	"CRM/go/authService/internal/logger"
	"CRM/go/authService/internal/proto/authService"
	"CRM/go/authService/pkg/hash"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

func Authorization(request *authService.AuthorizationRequest, response *authService.AuthorizationResponse) {
	request.Account.Password = hash.GenerateHash(request.Account.Password + config.GetConfig().Secret)

	row := postgres.GetAccountByLogin(request.Account)

	var lastActivity time.Time
	var dateCreated time.Time

	err := row.Scan(&request.Account.Id, &lastActivity, &dateCreated)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan account: %v", err), "accountLogin", request.Account.Login)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "account not found",
			HttpStatus:   int64(http.StatusNotFound),
		}
		return
	}

	request.Account.LastActivity = timestamppb.New(lastActivity)
	request.Account.DateCreated = timestamppb.New(dateCreated)

	request.Account.LastActivity = timestamppb.New(time.Now().UTC())

	_, err = postgres.UpdateLastActivityByAccountId(request.Account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("update last activity: %v", err), "accountId", request.Account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error update last activity",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	err = DeleteAllSessionsByAccount(request.Account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("delete all sessions by account: %v", err), "accountId", request.Account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error delete all sessions by account",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	response.Session, err = CreateSession(request.Account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("create session: %v", err), "accountId", request.Account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error create session",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	logger.CreateLog("info", "authorization successfully", "accountId", request.Account.Id)
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "authorization successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func CreateAccount(request *authService.RegistrationRequest, response *authService.RegistrationResponse) {
	response.Account = &authService.Account{
		Login:        request.Account.Login,
		Password:     hash.GenerateHash(request.Account.Password + config.GetConfig().Secret),
		LastActivity: timestamppb.New(time.Now().UTC()),
		DateCreated:  timestamppb.New(time.Now().UTC()),
	}

	row := postgres.CreateAccount(response.Account) //Можно обработать повторение логина
	err := row.Scan(&response.Account.Id)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("create account: %v", err))
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error create account",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	response.Session, err = CreateSession(response.Account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("create session: %v", err), "accountId", response.Account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error create session",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	logger.CreateLog("info", "create account successfully", "accountId", response.Account.Id)
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "create account successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func Logout(request *authService.LogoutRequest, response *authService.LogoutResponse) {
	_, err := postgres.DeleteSessionByAccountId(request.Session)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("delete session by account id: %v", err))
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error delete session by account id",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	err = redis.Set(context.Background(), fmt.Sprintf("exp:%v", request.Session.AccessToken), "", time.Second).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("delete key from redis: %v", err))
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error delete key from redis",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	logger.CreateLog("info", "logout successfully")
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "logout successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func CheckAuthorization(request *authService.CheckAuthorizationRequest, response *authService.CheckAuthorizationResponse) {
	val, err := redis.Get(context.Background(), request.Session.AccessToken).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			logger.CreateLog("info", fmt.Sprintf("key not found in redis: %v", err))
		} else {
			logger.CreateLog("error", fmt.Sprintf("error get key from redis: %v", err))
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error get key from redis",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}
	} else {
		account := &authService.Account{}
		err = json.Unmarshal([]byte(val), &account)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("unmarshal account: %v", err))
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error unmarshal account",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}
		account.LastActivity = timestamppb.New(time.Now().UTC())
		account.Login = ""
		account.Password = ""

		data, err := json.Marshal(account)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("marshal account: %v", err))
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error marshal account",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}

		err = redis.Set(context.Background(), request.Session.AccessToken, data, 0).Err()
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("set key in redis: %v", err))
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error set key in redis",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}

		logger.CreateLog("info", "check authorization successfully in redis")
		response.Status = &authService.Status{
			Successfully: true,
			Message:      "check authorization successfully in redis",
			HttpStatus:   http.StatusOK,
		}
		return
	}

	var dateExpirationAccessToken time.Time

	row := postgres.GetSessionByAccessToken(request.Session)
	err = row.Scan(&request.Session.Id, &request.Session.IdAccount, &dateExpirationAccessToken)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan session: %v", err), "accessToken", request.Session.AccessToken)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error check authorization",
			HttpStatus:   int64(http.StatusUnauthorized),
		}
		return
	}

	request.Session.DateExpirationAccessToken = timestamppb.New(dateExpirationAccessToken)

	if request.Session.DateExpirationAccessToken.AsTime().UTC().Compare(time.Now().UTC()) < 0 {
		_, err = postgres.RemoveAccessTokenBySessionId(request.Session)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("remove access token by session id: %v", err), "accessToken", request.Session.AccessToken)
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error remove access token by session id",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}

		logger.CreateLog("info", "access token expired", "accessToken", request.Session.AccessToken)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "access token expired",
			HttpStatus:   int64(http.StatusUnauthorized),
		}
		return
	}

	account := &authService.Account{
		Id:           request.Session.IdAccount,
		LastActivity: timestamppb.New(time.Now().UTC()),
	}

	_, err = postgres.UpdateLastActivityByAccountId(account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("update last activity by account id: %v", err), "accountId", account)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error update last activity by account id",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	logger.CreateLog("info", "check authorization successfully in postgres")
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "check authorization successfully in postgres",
		HttpStatus:   http.StatusOK,
	}
	return
}

func UpdateAccessToken(request *authService.UpdateAccessTokenRequest, response *authService.UpdateAccessTokenResponse) {
	response.Session = &authService.Session{
		RefreshToken:              request.Session.RefreshToken,
		AccessToken:               hash.GenerateRandomHash(),
		DateExpirationAccessToken: timestamppb.New(time.Now().UTC().Add(time.Minute)),
	}

	row := postgres.GetSessionByRefreshToken(response.Session)

	var dateExpirationRefreshToken time.Time

	err := row.Scan(&response.Session.Id, &response.Session.IdAccount, &dateExpirationRefreshToken)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan session: %v", err), "refreshToken", response.Session.RefreshToken)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error check authorization",
			HttpStatus:   int64(http.StatusUnauthorized),
		}
		return
	}

	response.Session.DateExpirationRefreshToken = timestamppb.New(dateExpirationRefreshToken)

	if response.Session.DateExpirationRefreshToken.AsTime().UTC().Compare(time.Now().UTC()) < 0 {
		_, err = postgres.RemoveSessionBySessionId(response.Session)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("remove session by session id: %v", err), "refreshToken", response.Session.RefreshToken)
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error remove session by session id",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}

		logger.CreateLog("info", "refresh token expired", "refreshToken", response.Session.RefreshToken)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "refresh token expired",
			HttpStatus:   int64(http.StatusUnauthorized),
		}
		return
	}

	_, err = postgres.UpdateAccessTokenByRefreshToken(response.Session)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("update access token by refresh token: %v", err), "refreshToken", response.Session.RefreshToken)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error update access token by refresh token",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	account := &authService.Account{
		Id: response.Session.IdAccount,
	}

	var lastActivity time.Time
	var dateCreated time.Time

	row = postgres.GetAccountById(account)
	err = row.Scan(&account.Login, &account.Password, &lastActivity, &dateCreated)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan account: %v", err), "accountId", account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "account not found",
			HttpStatus:   int64(http.StatusNotFound),
		}
		return
	}

	account.Login = ""
	account.Password = ""
	account.LastActivity = timestamppb.New(lastActivity)
	account.DateCreated = timestamppb.New(dateCreated)

	data, err := json.Marshal(account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error marshal account: %v", err), "userId", account.Id)
	}

	err = redis.Set(context.Background(), response.Session.AccessToken, data, 0).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err), "userId", account.Id)
	}

	err = redis.Set(context.Background(), fmt.Sprintf("exp:%v", response.Session.AccessToken), "", time.Minute).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("error set key in redis: %v", err), "userId", account.Id)
	}

	logger.CreateLog("info", "update access token successfully")
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "update access token successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetAccount(request *authService.GetAccountByAccessTokenRequest, response *authService.GetAccountByAccessTokenResponse) {
	session := &authService.Session{
		AccessToken: request.Session.AccessToken,
	}

	val, err := redis.Get(context.Background(), session.AccessToken).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			logger.CreateLog("info", fmt.Sprintf("key not found in redis: %v", err))
		} else {
			logger.CreateLog("error", fmt.Sprintf("error get key from redis: %v", err))
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error get key from redis",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}
	} else {
		err = json.Unmarshal([]byte(val), &response.Account)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("error unmarshal account: %v", err))
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error unmarshal account",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}

		logger.CreateLog("info", "get account successfully in redis")
		response.Status = &authService.Status{
			Successfully: true,
			Message:      "get account successfully in redis",
			HttpStatus:   http.StatusOK,
		}
		return
	}

	var dateExpirationAccessToken time.Time

	row := postgres.GetSessionByAccessToken(session)
	err = row.Scan(&session.Id, &session.IdAccount, &dateExpirationAccessToken)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan session: %v", err), "accessToken", session.AccessToken)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error get account",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	session.DateExpirationAccessToken = timestamppb.New(dateExpirationAccessToken)

	response.Account = &authService.Account{
		Id: session.IdAccount,
	}

	var lastActivity time.Time
	var dateCreated time.Time

	row = postgres.GetAccountById(response.Account)
	err = row.Scan(&response.Account.Login, &response.Account.Password, &lastActivity, &dateCreated)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan account: %v", err), "accountId", response.Account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error get account",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	response.Account.Login = ""
	response.Account.Password = ""
	response.Account.LastActivity = timestamppb.New(lastActivity)
	response.Account.DateCreated = timestamppb.New(dateCreated)

	logger.CreateLog("info", "get account successfully in postgres")
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "get account successfully in postgres",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetAccountById(request *authService.GetAccountByIdRequest, response *authService.GetAccountByIdResponse) {
	response.Account = &authService.Account{
		Id: request.Account.Id,
	}

	var lastActivity time.Time
	var dateCreated time.Time

	row := postgres.GetAccountById(response.Account)
	err := row.Scan(&response.Account.Login, &response.Account.Password, &lastActivity, &dateCreated)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("scan account: %v", err), "accountId", response.Account.Id)
		response.Status = &authService.Status{
			Successfully: false,
			Message:      "error get account",
			HttpStatus:   int64(http.StatusInternalServerError),
		}
		return
	}

	response.Account.Login = ""
	response.Account.Password = ""
	response.Account.LastActivity = timestamppb.New(lastActivity)
	response.Account.DateCreated = timestamppb.New(dateCreated)

	logger.CreateLog("info", "get account successfully")
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "get account successfully",
		HttpStatus:   http.StatusOK,
	}
	return
}

func GetAccounts(request *authService.GetAccountsRequest, response *authService.GetAccountsResponse) {
	response.Accounts = make(map[int64]*authService.Account)
	for _, id := range request.Id {
		account := &authService.Account{
			Id: id,
		}
		var lastActivity time.Time
		var dateCreated time.Time

		row := postgres.GetAccountById(account)
		err := row.Scan(&account.Login, &account.Password, &lastActivity, &dateCreated)
		if err != nil {
			logger.CreateLog("error", fmt.Sprintf("scan account: %v", err), "accountId", account.Id)
			response.Status = &authService.Status{
				Successfully: false,
				Message:      "error get accounts",
				HttpStatus:   int64(http.StatusInternalServerError),
			}
			return
		}

		account.Login = ""
		account.Password = ""
		account.LastActivity = timestamppb.New(lastActivity)
		account.DateCreated = timestamppb.New(dateCreated)

		response.Accounts[account.Id] = account
	}

	logger.CreateLog("info", "get accounts successfully")
	response.Status = &authService.Status{
		Successfully: true,
		Message:      "get accounts successfully",
		HttpStatus:   http.StatusOK,
	}
}

func DeleteAllSessionsByAccount(account *authService.Account) error {
	_, err := postgres.DeleteAllSessionsByAccount(account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("delete all sessions by account: %v", err), "accountId", account.Id)
		return err
	}

	return nil
}

func CreateSession(account *authService.Account) (*authService.Session, error) {
	session := &authService.Session{
		IdAccount:                  account.Id,
		AccessToken:                hash.GenerateRandomHash(),
		DateExpirationAccessToken:  timestamppb.New(time.Now().UTC().Add(time.Minute)),
		RefreshToken:               hash.GenerateRandomHash(),
		DateExpirationRefreshToken: timestamppb.New(time.Now().UTC().Add(time.Minute * 5)),
	}

	_, err := postgres.CreateSession(session)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("create session: %v", err), "accountId", account.Id)
		return nil, err
	}

	account.Login = ""
	account.Password = ""

	data, err := json.Marshal(account)
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("marshal account: %v", err), "accountId", account.Id)
		return nil, err
	}

	err = redis.Set(context.Background(), session.AccessToken, data, 0).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("set key in redis: %v", err), "accountId", account.Id)
		return nil, err
	}

	err = redis.Set(context.Background(), fmt.Sprintf("exp:%v", session.AccessToken), "", time.Minute).Err()
	if err != nil {
		logger.CreateLog("error", fmt.Sprintf("set key in redis: %v", err), "accountId", account.Id)
		return nil, err
	}

	return session, nil
}
