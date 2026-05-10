package routes

import (
	"fileserver/internal/auth/middleware"
	"fileserver/internal/filesystem/handler"

	"github.com/gin-gonic/gin"
)

func RegisterFilesystemRoutes(
	rg *gin.RouterGroup,
	handler *handler.FilesystemHandler,
) {

	fs := rg.Group("/fs")

	fs.Use(middleware.JWTMiddleware())

	{
		fs.GET("/list", handler.List)

		fs.GET("/stat", handler.Stat)

		fs.POST("/mkdir", handler.Mkdir)

		fs.DELETE("/delete", handler.Delete)

		fs.POST("/move", handler.Move)

		fs.POST("/copy", handler.Copy)
	}
}
