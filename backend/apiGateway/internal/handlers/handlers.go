package handlers

import (
	"CRM/go/apiGateway/internal/api"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func Authorization(c *gin.Context) {

	resp, err := api.NewRequest(http.MethodPost, "http://localhost:3001/api/auth", c.Request)

	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"successfully": false,
			"error":        "error creating request",
		})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"successfully": false,
			"error":        "error reading response body",
		})
		return
	}

	var bodyJSON map[string]any
	json.Unmarshal(body, &bodyJSON)

	_, ok := bodyJSON["successfully"]
	if !ok {
		fmt.Printf("error response body: %v\n", bodyJSON)
		c.JSON(http.StatusInternalServerError, gin.H{
			"successfully": false,
			"error":        "error response body",
		})
		return
	}

	if !bodyJSON["successfully"].(bool) {
		fmt.Printf("error response body: %v\n", bodyJSON)
		c.JSON(http.StatusInternalServerError, bodyJSON)
		return
	}

	c.JSON(http.StatusOK, bodyJSON)
}

func Registration(c *gin.Context) {

	resp, err := api.NewRequest(http.MethodPost, "http://localhost:3001/api/reg", c.Request)

	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"successfully": false,
			"error":        "error creating request",
		})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"successfully": false,
			"error":        "error reading response body",
		})
		return
	}

	var bodyJSON map[string]any
	json.Unmarshal(body, &bodyJSON)
	fmt.Println(bodyJSON)

	_, ok := bodyJSON["successfully"]
	if !ok {
		fmt.Printf("error response body: %v\n", bodyJSON)
		c.JSON(http.StatusInternalServerError, gin.H{
			"successfully": false,
			"error":        "error response body",
		})
		return
	}

	if !bodyJSON["successfully"].(bool) {
		fmt.Printf("error response body: %v\n", bodyJSON)
		c.JSON(http.StatusInternalServerError, bodyJSON)
		return
	}

	c.JSON(http.StatusOK, bodyJSON)
}

func CheckAuthorization(c *gin.Context) {

	resp, err := api.NewRequest(http.MethodPost, "http://localhost:3001/api/checkAuth", c.Request)

	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"successfully": false,
			"error":        "error creating request",
		})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"successfully": false,
			"error":        "error reading response body",
		})
		return
	}

	var bodyJSON map[string]any
	json.Unmarshal(body, &bodyJSON)
	fmt.Println(bodyJSON)

	_, ok := bodyJSON["successfully"]
	if !ok {
		fmt.Printf("error response body: %v\n", bodyJSON)
		c.JSON(http.StatusOK, gin.H{
			"successfully": false,
			"error":        "error response body",
		})
		return
	}

	/*if !bodyJSON["successfully"].(bool) {
		fmt.Printf("error response body: %v\n", bodyJSON)
		c.JSON(http.StatusOK, bodyJSON)
		return
	}*/

	c.JSON(http.StatusOK, bodyJSON)
}
