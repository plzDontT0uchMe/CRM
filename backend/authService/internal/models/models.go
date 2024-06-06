package models

import (
	"time"
)

type Account struct {
	Id           int       `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:"password"`
	Role         int       `json:"role"`
	LastActivity time.Time `json:"last_activity"`
	DateCreated  time.Time `json:"date_created"`
}

type Session struct {
	Id                         int       `json:"id"`
	IdAccount                  int       `json:"id_account"`
	AccessToken                string    `json:"access_token"`
	DateExpirationAccessToken  time.Time `json:"date_expiration_access_token"`
	RefreshToken               string    `json:"refresh_token"`
	DateExpirationRefreshToken time.Time `json:"date_expiration_refresh_token"`
}
