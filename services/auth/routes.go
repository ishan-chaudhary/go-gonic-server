package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/signup", Signup)
	router.POST("/login", Login)
}
