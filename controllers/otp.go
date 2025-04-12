package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/models"
	"github.com/FLXA-Auth-Service-IOTA/utilities"
	"github.com/gin-gonic/gin"
)

type charsetType string

const (
	Alphanumeric charsetType = "alphanumeric"
	Numeric      charsetType = "numeric"
)

func VerifyOTP(c *gin.Context) {
	ctx := context.Background()
	var user models.User

	userId := c.Query("userId")
	otp := c.Query("otp")
	phoneNumber := fmt.Sprintf("+%s", c.Query("phoneNumber"))
	OtpKeyName := "otp_code:" + phoneNumber
	val, err := initializers.RDB.Get(ctx, OtpKeyName).Result()

	if err != nil {
		c.JSON(404, gin.H{
			"message": "No key found",
		})
		return
	}

	if val != otp {
		c.JSON(401, gin.H{
			"message": "Unauthorized: Wrong OTP",
		})
		return
	}

	userTransactionErr := initializers.DB.Raw("SELECT * FROM \"User\" WHERE \"User\".user_id = ?", userId).Scan(&user).Error
	message, code, details := utilities.TransactionErrorHandler(userTransactionErr)
	if code != 0 {
		c.JSON(code, gin.H{
			"message": message,
			"details": details,
		})
		return
	}
	if user.UserID == 0 {
		c.JSON(404, gin.H{
			"message": "User not found",
		})
		return
	}

	result := initializers.DB.Exec("UPDATE \"Card\" SET card_status = 'VERIFIED' WHERE card_phone_number = ?", phoneNumber)
	transactionError := result.Error
	message, code, details = utilities.TransactionErrorHandler(transactionError)
	if code != 0 {
		c.JSON(code, gin.H{
			"message": message,
			"details": details,
		})
		return
	}

	SessionKeyName := user.Name + strconv.Itoa(int(user.UserID))
	hashedKeyName, hashErr := utilities.Hash(SessionKeyName)
	if hashErr != nil {
		c.JSON(500, gin.H{
			"message": hashErr.Error(),
		})
	}
	sessionID := "session_tokens:" + hashedKeyName
	initializers.RDB.Set(ctx, sessionID, user.UserID, 8*time.Hour)

	initializers.RDB.Set(ctx, "otp_code:"+phoneNumber, "", -1*time.Second)

	c.JSON(200, map[string]interface{}{
		"message":     "Phone number verified",
		"phoneNumber": phoneNumber,
		"session":     hashedKeyName,
	})
}

func ResendOTP(c *gin.Context) {
	phoneNumber := c.Query("phoneNumber")
	otpCode, _ := generateRandomString(6, Numeric)
	ctx := context.Background()
	sessionID := "otp_code:+" + phoneNumber
	initializers.RDB.Set(ctx, sessionID, otpCode, 3*time.Minute)

	OtpServiceUrl := fmt.Sprintf("%s/send-otp?number=+%s&otp=%s", os.Getenv("FLXA_OTP_SERVICE"), strings.Trim(phoneNumber, " "), otpCode)
	fmt.Println(phoneNumber, otpCode)
	fmt.Println(sessionID)
	response, err := http.Get(OtpServiceUrl)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "An error occurred on Message Service",
		})
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "An error occurred on Message Service",
		})
	}
	fmt.Println(string(responseData))
	resp := "Resent OTP to " + phoneNumber
	c.JSON(200, gin.H{
		"message": resp,
	})
}
