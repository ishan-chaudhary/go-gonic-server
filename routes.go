package main

import (
	"swiggy/gin/services/auth"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes(router *gin.Engine) {
	auth.AuthRoutes(router)
}
