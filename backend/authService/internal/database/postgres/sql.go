package postgres

import (
	"CRM/go/authService/internal/model"
	"fmt"
	"net/http"
)

func GetUser(user *model.User) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE login = $1 AND password = $2 LIMIT 1", user.Login, user.Password)

	err = row.Scan(&user.Id)
	if err != nil {
		fmt.Println("error binding JSON for user: ", err)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK

}

func CreateUser(user *model.User) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", user.Login, user.Password)

	var tempUser model.User
	err = row.Scan(&tempUser.Id)
	if err != nil {
		fmt.Println("error binding JSON for tempUser: ", err)
		return err, http.StatusOK
	}

	user.Id = tempUser.Id

	return nil, http.StatusOK
}

func RemoveAllSessionsByUser(user *model.User) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE id_user = $1", user.Id)

	if err != nil {
		fmt.Println("database query error: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK

}

func CreateSession(session *model.Session) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO sessions (id_user, access_token, date_expiration_access_token, refresh_token, date_expiration_refresh_token) VALUES ($1, $2, $3, $4, $5)",
		session.UserId, session.AccessToken, session.DateExpirationAccessToken, session.RefreshToken, session.DateExpirationRefreshToken)

	if err != nil {
		fmt.Println("database query error: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK
}

func GetSessionByAccessToken(session *model.Session) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, date_expiration_access_token FROM sessions WHERE access_token = $1", session.AccessToken)

	err = row.Scan(&session.Id, &session.DateExpirationAccessToken)
	if err != nil {
		fmt.Println("error binding JSON for session: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK
}

func RemoveAccessTokenBySessionId(session *model.Session) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	_, err = db.Exec("UPDATE sessions SET access_token = NULL, date_expiration_access_token = NULL WHERE id = $1", session.Id)

	if err != nil {
		fmt.Println("database query error: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK
}

func GetSessionByRefreshToken(session *model.Session) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, date_expiration_refresh_token FROM sessions WHERE refresh_token = $1", session.RefreshToken)

	err = row.Scan(&session.Id, &session.DateExpirationRefreshToken)
	if err != nil {
		fmt.Println("error binding JSON for session: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK
}

func RemoveSessionBySessionId(session *model.Session) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE id = $1", session.Id)

	if err != nil {
		fmt.Println("database query error: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK
}

func UpdateAccessTokenByRefreshToken(session *model.Session) (error, int) {

	db, err := Connect()
	if err != nil {
		fmt.Println("database connection error: ", err)
		return err, http.StatusOK
	}
	defer db.Close()

	_, err = db.Exec("UPDATE sessions SET access_token = $1, date_expiration_access_token = $2 WHERE refresh_token = $3",
		session.AccessToken, session.DateExpirationAccessToken, session.RefreshToken)

	if err != nil {
		fmt.Println("database query error: ", err)
		return err, http.StatusOK
	}

	return nil, http.StatusOK
}
