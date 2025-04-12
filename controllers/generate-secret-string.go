package controllers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/models"
	"github.com/FLXA-Auth-Service-IOTA/utilities"
	"github.com/gin-gonic/gin"
)

type generateSecretStringBody struct {
	UserID      int
	PhoneNumber string
}

func generateRandomString(length int, charsetType charsetType) (string, error) {
	const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const numericCharset = "0123456789"

	var charset string
	switch charsetType {
	case Alphanumeric:
		charset = alphanumericCharset
	case Numeric:
		charset = numericCharset
	default:
		return "", errors.New("invalid charset type; must be 'alphanumeric' or 'numeric'")
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))] // Use the local generator's Intn
	}

	return string(result), nil
}

func validateGenerateSecretStringBody(body *generateSecretStringBody) (string, int, []string) {
	message := ""
	code := 0
	var details []string

	if body.PhoneNumber == "" {
		message = "Invalid Body"
		code = 400
		details = append(details, "Phone Number is required")
	}

	return message, code, details
}

func GenerateSecretString(c *gin.Context) {
	var body generateSecretStringBody
	var user models.User

	jsonErr := c.ShouldBindJSON(&body)
	message, code, details := utilities.ValidateJSON(jsonErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	secretString, generatingErr := generateRandomString(8, Alphanumeric)
	hashedSecretString, hashErr := utilities.Hash(secretString)
	if generatingErr != nil {
		fmt.Println(generatingErr.Error())
		c.JSON(500, gin.H{
			"message": generatingErr.Error(),
		})
		return
	}
	if hashErr != nil {
		c.JSON(500, gin.H{
			"message": hashErr.Error(),
		})
		return
	}

	message, code, details = validateGenerateSecretStringBody(&body)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	query := "SELECT user_id FROM \"Card\" WHERE card_phone_number = ?"
	result := initializers.DB.Raw(query, "+"+body.PhoneNumber).Scan(&user)
	fmt.Println(user)
	fmt.Println(body.PhoneNumber)
	transactionErr := result.Error
	message, code, details = utilities.TransactionErrorHandler(transactionErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	query = "UPDATE \"User\" SET secret_string = ? WHERE \"User\".user_id = ?;"
	result = initializers.DB.Exec(query, hashedSecretString, user.UserID)
	transactionErr = result.Error
	message, code, details = utilities.TransactionErrorHandler(transactionErr)
	if code != 0 {
		c.JSON(code, map[string]interface{}{
			"message": message,
			"details": details,
		})
		return
	}

	message = fmt.Sprintf("Sent Secret String: %s, to %s", secretString, body.PhoneNumber)

	OtpServiceUrl := fmt.Sprintf("%s/send-message?number=+%s&message=%s", os.Getenv("FLXA_OTP_SERVICE"), body.PhoneNumber, secretString)
	fmt.Println("LOG " + OtpServiceUrl)
	response, err := http.Get(OtpServiceUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
	c.JSON(200, gin.H{
		"message": message,
	})
}
