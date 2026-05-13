package routes

import (
	authmiddleware "fileserver/internal/auth/middleware"
	"fileserver/internal/torrent/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTorrentRoutes(
	rg *gin.RouterGroup,
	handler *handler.TorrentHandler,
) {

	torrent := rg.Group("/torrent")

	torrent.Use(
		authmiddleware.JWTMiddleware(),
	)

	{
		torrent.POST(
			"/add",
			handler.AddTorrent,
		)

		torrent.DELETE(
			"/:id",
			handler.RemoveTorrent,
		)
	}
}
