package controllers

import (
	"fmt"

	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func CheckSession(c *gin.Context) {
	ctx := context.Background()
	token := c.Query("token")

	keyName := "session_tokens:" + token
	fmt.Println(keyName)

	val, err := initializers.RDB.Get(ctx, keyName).Result()
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Token Not Found or Expired",
		})
		return
	}

	c.JSON(200, gin.H{
		"value": val,
	})
}
