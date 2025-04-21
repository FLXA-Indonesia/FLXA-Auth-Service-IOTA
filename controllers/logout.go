package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/utilities"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	ctx := context.Background()
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{
			"message": "token is required",
		})
		return
	}

	err := initializers.RDB.Set(ctx, "session_tokens:"+token, "", 1*time.Second).Err()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Log out success",
	})
}

type LogoutAllBody struct {
	UserId int
	Name   string
}

func LogoutAll(c *gin.Context) {
	var body LogoutAllBody
	ctx := context.Background()

	jsonErr := c.ShouldBindJSON(&body)

	fmt.Println(body)

	keyName := body.Name + strconv.Itoa(body.UserId)
	hashedKeyName, hashErr := utilities.Hash(keyName)
	if hashErr != nil {
		c.JSON(500, gin.H{
			"message": hashErr.Error(),
		})
	}
	sessionID := "session_tokens:" + hashedKeyName

	if jsonErr != nil {
		c.JSON(415, gin.H{
			"message": "Invalid JSON in body",
		})
		return
	}

	initializers.RDB.Set(ctx, sessionID, "", -1*time.Second)

	c.JSON(200, gin.H{
		"message": "Session " + sessionID + " deleted",
		"data":    body,
	})
}
