package handler

import (
	"fileserver/internal/jobs/service"
	"net/http"

	apperrors "fileserver/pkg/errors"
	"fileserver/pkg/response"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	Service *service.JobService
}

func NewJobHandler(
	service *service.JobService,
) *JobHandler {

	return &JobHandler{
		Service: service,
	}
}

func (h *JobHandler) GetJob(
	c *gin.Context,
) {

	jobID := c.Param("id")

	job, exists := h.Service.GetJob(jobID)

	if !exists {

		apperrors.Respond(
			c,
			apperrors.ErrNotFound,
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		job,
	)
}

func (h *JobHandler) ListJobs(
	c *gin.Context,
) {

	response.Success(
		c,
		http.StatusOK,
		h.Service.ListJobs(),
	)
}
