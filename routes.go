package main

import (
	"lingjiao0710/ginEssential/controller"
	"lingjiao0710/ginEssential/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("api/auth/login", controller.Login)

	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)

	return r
}
