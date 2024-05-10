package main

import (
	"CRM/go/authService/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func getUser(c *gin.Context) {

	dbHandler, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.Id, &user.Login, &user.Password); err != nil {
			log.Fatal(err)
			return
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(users)
	c.JSON(http.StatusOK, users)
}

func main() {
	server := gin.Default()
	server.GET("/test", getUser)
	server.Run()
}
