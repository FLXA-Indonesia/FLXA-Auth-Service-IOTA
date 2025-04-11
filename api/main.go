package handler

import (
	"net/http"

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

	rGin.POST("/register")
	rGin.POST("/login")
	rGin.POST("/generate-secret")
	rGin.GET("/verify-otp")
	rGin.GET("/check-session")
	rGin.GET("/logout")
	rGin.PATCH("/complete-profile")
	rGin.GET("/resend-otp")
	rGin.POST("/logout-all")

	rGin.ServeHTTP(w, r)
}
