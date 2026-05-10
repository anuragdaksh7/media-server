package routes

import (
	authmiddleware "fileserver/internal/auth/middleware"
	"fileserver/internal/jobs/handler"

	"github.com/gin-gonic/gin"
)

func RegisterJobRoutes(
	rg *gin.RouterGroup,
	handler *handler.JobHandler,
) {

	jobs := rg.Group("/jobs")

	jobs.Use(
		authmiddleware.JWTMiddleware(),
	)

	{
		jobs.GET("", handler.ListJobs)

		jobs.GET("/:id", handler.GetJob)
	}
}
