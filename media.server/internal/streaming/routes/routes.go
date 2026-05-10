package routes

import (
	"fileserver/internal/streaming/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStreamingRoutes(
	rg *gin.RouterGroup,
	handler *handler.StreamingHandler,
) {

	stream := rg.Group("stream")

	{
		stream.GET(
			"/video",
			handler.Stream,
		)

		stream.GET(
			"/download",
			handler.Download,
		)
	}
}
