package routes

import (
	authmiddleware "fileserver/internal/auth/middleware"
	"fileserver/internal/upload/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(
	rg *gin.RouterGroup,
	handler *handler.UploadHandler,
) {

	upload := rg.Group("/upload")

	upload.Use(
		authmiddleware.JWTMiddleware(),
	)

	{
		upload.POST(
			"/file",
			handler.Upload,
		)
	}
}
