package models

import "time"

type User struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Role       int    `json:"role"`
}

type Session struct {
	Id                         int       `json:"id"`
	UserId                     int       `json:"id_user"`
	AccessToken                string    `json:"access_token"`
	DateExpirationAccessToken  time.Time `json:"date_expiration_access_token"`
	RefreshToken               string    `json:"refresh_token"`
	DateExpirationRefreshToken time.Time `json:"date_expiration_refresh_token"`
}
