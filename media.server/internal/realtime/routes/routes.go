package routes

import (
	authmiddleware "fileserver/internal/auth/middleware"
	"fileserver/internal/realtime/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRealtimeRoutes(
	rg *gin.RouterGroup,
	handler *handler.RealtimeHandler,
) {

	realtime := rg.Group("/realtime")

	realtime.Use(
		authmiddleware.JWTMiddleware(),
	)

	{
		realtime.GET(
			"/ws",
			handler.WebSocket,
		)
	}
}
