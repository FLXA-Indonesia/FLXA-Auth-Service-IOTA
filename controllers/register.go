package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/models"
	"github.com/FLXA-Auth-Service-IOTA/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerBody struct {
	PhoneNumber string
	Name        string
	Provider    string
}

func validateRegisterBody(body registerBody) (string, int, []string) {
	message := ""
	code := 0
	var detail []string

	if !strings.HasPrefix(body.PhoneNumber, "+62") {
		message = "Invalid phone number"
		code = 400
		detail = append(detail, "Phone number must start with country code (exclusive +62 for now)")
	}

	if len(body.PhoneNumber) < 10 {
		message = "Invalid phone number"
		code = 400
		detail = append(detail, "Phone number length is invalid")
	}

	if body.Name == "" {
		message = "Invalid request body"
		code = 400
		detail = append(detail, "Name is required")
	}

	if body.Provider == "" {
		message = "Invalid request body"
		code = 400
		detail = append(detail, "Provider is required")
	}

	return message, code, detail
}

func Register(c *gin.Context) {
	var body registerBody
	var newCard models.Card
	var newUser models.User

	jsonErr := c.ShouldBindJSON(&body)
	message, code, details := utilities.ValidateJSON(jsonErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}
	fmt.Println(body)

	message, code, details = validateRegisterBody(body)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	tx := initializers.DB.Begin()

	newUser.Name = body.Name
	userTransactionResult := tx.Create(&newUser)
	userTransactionErr := userTransactionResult.Error
	message, code, details = utilities.TransactionErrorHandler(userTransactionErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		tx.Rollback()
		return
	}

	newCard.CardPhoneNumber = body.PhoneNumber
	operatorMapping := map[string]string{
		"Telkomsel": "58fbf693-3363-4066-9ac3-489288d950c9",
		"Indosat":   "c9b0f2b1-f9e4-4105-9e61-614024fbcdb3",
		"XL Axiata": "0dbfa48b-4d97-4ce3-acbc-84910b911e1e",
	}
	operatorUUID, _ := uuid.Parse(operatorMapping[body.Provider])
	newCard.OperatorID = operatorUUID
	newCard.UserID = newUser.UserID
	newCard.CardStatus = "NOT VERIFIED"
	newCard.CardDateAdded = time.Now()
	phoneTransactionResult := tx.Create(&newCard)
	phoneTransactionErr := phoneTransactionResult.Error
	message, code, details = utilities.TransactionErrorHandler(phoneTransactionErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		tx.Rollback()
		return
	}

	balanceTransactionErr := tx.Exec("INSERT INTO \"Balance\" (user_id, balance_amount) VALUES ($1, 0)", newUser.UserID).Error
	message, code, details = utilities.TransactionErrorHandler(balanceTransactionErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		tx.Rollback()
		return
	}

	otpCode, _ := generateRandomString(6, Numeric)
	ctx := context.Background()
	sessionID := "otp_code:" + newCard.CardPhoneNumber
	initializers.RDB.Set(ctx, sessionID, otpCode, 3*time.Minute)

	OtpServiceUrl := fmt.Sprintf("%s/send-otp?number=%s&otp=%s", os.Getenv("FLXA_OTP_SERVICE"), body.PhoneNumber, otpCode)
	response, err := http.Get(OtpServiceUrl)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "An error occurred on Message Service",
		})
		tx.Rollback()
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "An error occurred on Message Service",
		})
		tx.Rollback()
	}
	fmt.Println(string(responseData))

	tx.Commit()
	c.JSON(201, map[string]interface{}{
		"message":     "Register successful",
		"id":          newUser.UserID,
		"phoneNumber": newCard.CardPhoneNumber,
		"name":        newUser.Name,
	})
}
