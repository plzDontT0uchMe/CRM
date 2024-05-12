package handlers

import (
	"CRM/go/authService/internal/model"
	"CRM/go/authService/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorization(c *gin.Context) {

	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		fmt.Printf("error binding json for user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"successfully": false,
			"error":        "error binding json for user",
		})
		return
	}

	err, httpStatus := service.AuthorizeUser(&user)

	if err != nil {
		c.JSON(httpStatus, gin.H{
			"successfully": false,
			"error":        "error registering user",
		})
		return
	}

	err, httpStatus = service.RemoveAllSessionsByUser(&user)

	if err != nil {
		c.JSON(httpStatus, gin.H{
			"successfully": false,
			"error":        "error removing all sessions by user",
		})
		return
	}

	session, err, httpStatus := service.CreateSession(&user)

	if err != nil {
		c.JSON(httpStatus, gin.H{
			"successfully": false,
			"error":        "error creating session",
		})
		return
	}

	c.JSON(httpStatus, gin.H{
		"successfully": true,
		"message":      "user created successfully, session created successfully",
		"session":      session,
	})

}

func Registration(c *gin.Context) {

	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		fmt.Printf("error binding json for user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"successfully": false,
			"error":        "error binding json for user",
		})
		return
	}

	err, httpStatus := service.RegisterUser(&user)

	if err != nil {
		c.JSON(httpStatus, gin.H{
			"successfully": false,
			"error":        "error registering user",
		})
		return
	}

	err, httpStatus = service.RemoveAllSessionsByUser(&user)

	if err != nil {
		c.JSON(httpStatus, gin.H{
			"successfully": false,
			"error":        "error removing all sessions by user",
		})
		return
	}

	session, err, httpStatus := service.CreateSession(&user)

	if err != nil {
		c.JSON(httpStatus, gin.H{
			"successfully": false,
			"error":        "error creating session",
		})
		return
	}

	c.JSON(httpStatus, gin.H{
		"successfully": true,
		"message":      "user created successfully, session created successfully",
		"session":      session,
	})

}

func CheckAuthorization(c *gin.Context) {

	var session model.Session
	err := c.BindJSON(&session)
	if err != nil {
		fmt.Printf("error binding json for session: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"successfully": false,
			"error":        "error binding json for session",
		})
		return
	}

	if session.AccessToken == "" && session.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"successfully": false,
			"error":        "access token and refresh token are empty",
		})
		return
	}

	m, flag, err, httpStatus := service.CheckAuthorization(&session)

	if err != nil {
		if session.AccessToken != "" {
			c.JSON(httpStatus, gin.H{
				"successfully": true,
				"error":        "error checking authorization by access token",
				"flag":         "getRefreshToken",
			})
			return
		}
		if session.RefreshToken != "" {
			c.JSON(httpStatus, gin.H{
				"successfully": false,
				"error":        "error checking authorization by refresh token",
				"flag":         "authorizationFailed",
			})
			return
		}
		return
	}

	c.JSON(httpStatus, gin.H{
		"successfully": true,
		"message":      m,
		"flag":         flag,
	})

}
