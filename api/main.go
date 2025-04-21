package handler

import (
	"net/http"

	"github.com/FLXA-Auth-Service-IOTA/controllers"
	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.ConnectRedis()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	rGin := gin.Default()

	rGin.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to FLXA Auth Service",
		})
	})

	rGin.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "pong",
		})
	})

	rGin.POST("/register", controllers.Register)
	rGin.POST("/login", controllers.Login)
	rGin.POST("/generate-secret", controllers.GenerateSecretString)
	rGin.GET("/verify-otp", controllers.VerifyOTP)
	rGin.GET("/check-session", controllers.CheckSession)
	rGin.GET("/logout", controllers.Logout)
	rGin.PATCH("/complete-profile", controllers.CompleteProfile)
	rGin.GET("/resend-otp", controllers.ResendOTP)
	rGin.POST("/logout-all", controllers.LogoutAll)

	rGin.ServeHTTP(w, r)
}
