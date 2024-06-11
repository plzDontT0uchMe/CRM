package postgres

import (
	"CRM/go/authService/internal/proto/authService"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetAccountByLogin(account *authService.Account) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT id, last_activity, date_created FROM accounts WHERE login = $1 AND password = $2 LIMIT 1", account.Login, account.Password)
}

func GetAccountById(account *authService.Account) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT login, password, last_activity, date_created FROM accounts WHERE id = $1 LIMIT 1", account.Id)
}

func CreateAccount(account *authService.Account) pgx.Row {
	return GetDB().QueryRow(context.Background(), "INSERT INTO accounts (login, password, last_activity, date_created) VALUES ($1, $2, $3, $4) RETURNING id", account.Login, account.Password, account.LastActivity.AsTime(), account.DateCreated.AsTime())
}

func DeleteSessionByAccountId(session *authService.Session) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "DELETE FROM sessions WHERE access_token = $1", session.AccessToken)
}

func DeleteAllSessionsByAccount(account *authService.Account) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "DELETE FROM sessions WHERE id_account = $1", account.Id)
}

func CreateSession(session *authService.Session) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "INSERT INTO sessions (id_account, access_token, date_expiration_access_token, refresh_token, date_expiration_refresh_token) VALUES ($1, $2, $3, $4, $5)", session.IdAccount, session.AccessToken, session.DateExpirationAccessToken.AsTime(), session.RefreshToken, session.DateExpirationRefreshToken.AsTime())
}

func GetSessionByAccessToken(session *authService.Session) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT id, id_account, date_expiration_access_token FROM sessions WHERE access_token = $1 LIMIT 1", session.AccessToken)
}

func RemoveAccessTokenBySessionId(session *authService.Session) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE sessions SET access_token = NULL, date_expiration_access_token = NULL WHERE id = $1", session.Id)
}

func GetSessionByRefreshToken(session *authService.Session) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT id, id_account, date_expiration_refresh_token FROM sessions WHERE refresh_token = $1 LIMIT 1", session.RefreshToken)
}

func RemoveSessionBySessionId(session *authService.Session) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "DELETE FROM sessions WHERE id = $1", session.Id)
}

func UpdateAccessTokenByRefreshToken(session *authService.Session) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE sessions SET access_token = $1, date_expiration_access_token = $2 WHERE refresh_token = $3", session.AccessToken, session.DateExpirationAccessToken.AsTime(), session.RefreshToken)
}

func UpdateLastActivityByAccountId(account *authService.Account) (pgconn.CommandTag, error) {
	return GetDB().Exec(context.Background(), "UPDATE accounts SET last_activity = $1 WHERE id = $2", account.LastActivity.AsTime(), account.Id)
}
