package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/models"
	"github.com/FLXA-Auth-Service-IOTA/utilities"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginBody struct {
	PhoneNumber  string
	SecretString string
}

func validateLoginBody(body *loginBody) (string, int, []string) {
	message := ""
	code := 0
	var details []string
	if body.PhoneNumber == "" {
		message = "Invalid Body"
		code = 400
		details = append(details, "Phone number is required")
	}
	if body.SecretString == "" {
		message = "Invalid Body"
		code = 400
		details = append(details, "Secret string is required")
	}
	return message, code, details
}

func Login(c *gin.Context) {
	var user models.User
	var card models.Card
	var body loginBody
	ctx := context.Background()

	jsonErr := c.ShouldBindJSON(&body)
	message, code, details := utilities.ValidateJSON(jsonErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	message, code, details = validateLoginBody(&body)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	cardQuery := "SELECT * FROM \"Card\" WHERE card_phone_number = ?;"
	result := initializers.DB.Raw(cardQuery, body.PhoneNumber).Scan(&card)
	transactionErr := result.Error
	message, code, details = utilities.TransactionErrorHandler(transactionErr)

	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	cardExist := card.CardID != 0
	if !cardExist {
		details = append(details, "Phone number is not registered to any user")
		c.JSON(404, map[string]interface{}{
			"message": "Phone Number Not Found",
			"details": details,
		})
		return
	}

	userQuery := "SELECT * FROM \"User\" WHERE \"User\".user_id = ?;"
	result = initializers.DB.Raw(userQuery, card.UserID).Scan(&user)
	transactionErr = result.Error
	message, code, details = utilities.TransactionErrorHandler(transactionErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	areCredentialsValid := false
	comparationErr := bcrypt.CompareHashAndPassword([]byte(user.SecretString), []byte(body.SecretString))
	if comparationErr == nil {
		areCredentialsValid = true
	}
	fmt.Println(user.SecretString)
	fmt.Println(body.SecretString)
	fmt.Println(comparationErr)
	if !areCredentialsValid {
		details = append(details, "Invalid Login Credentials")
		c.JSON(401, map[string]interface{}{
			"message": "Invalid Credential Combination",
			"details": details,
			"errors":  comparationErr.Error(),
		})
		return
	}

	keyName := user.Name + strconv.Itoa(int(user.UserID))
	hashedKeyName, hashErr := utilities.Hash(keyName)
	if hashErr != nil {
		c.JSON(500, gin.H{
			"message": hashErr.Error(),
		})
	}
	sessionID := "session_tokens:" + hashedKeyName
	status := initializers.RDB.Set(ctx, sessionID, user.UserID, 8*time.Hour)
	fmt.Println(status.Result())
	c.JSON(200, map[string]interface{}{
		"message": "Login successful",
		"session": hashedKeyName,
		"user": map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
			"id":    user.UserID,
		},
		"expiresAt": time.Now().Add(8 * time.Hour),
	})
}
