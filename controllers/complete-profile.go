package controllers

import (
	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/utilities"
	"github.com/gin-gonic/gin"
)

type CompleteProfileBody struct {
	ProfilePhoto string
	Email        string
}

func validateCompleteProfileBody(body CompleteProfileBody) (string, int, []string) {
	message := ""
	code := 0
	var details []string

	if body.Email == "" {
		code = 400
		details = append(details, "Email is required")
	}

	// if body.ProfilePhoto == "" {
	// 	code = 400
	// 	details = append(details, "Profile picture is required")
	// }
	return message, code, details
}

func CompleteProfile(c *gin.Context) {
	var body CompleteProfileBody
	userId := c.Query("userId")

	jsonErr := c.ShouldBindJSON(&body)
	message, code, details := utilities.ValidateJSON(jsonErr)
	if code != 0 {
		c.JSON(code, gin.H{
			"messages": message,
			"details":  details,
		})
		return
	}

	message, code, details = validateCompleteProfileBody(body)
	if code != 0 {
		c.JSON(code, gin.H{
			"messages": message,
			"details":  details,
		})
		return
	}

	query := "UPDATE \"User\" SET email = ? WHERE user_id = ? "
	transactionErr := initializers.DB.Exec(query, body.Email, userId).Error
	message, code, details = utilities.TransactionErrorHandler(transactionErr)
	if code != 0 {
		c.JSON(code, gin.H{
			"messages": message,
			"details":  details,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Updated",
	})
}
