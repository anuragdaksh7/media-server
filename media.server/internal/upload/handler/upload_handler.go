package handler

import (
	"fileserver/internal/upload/dto"
	"fileserver/internal/upload/service"
	"net/http"

	apperrors "fileserver/pkg/errors"
	"fileserver/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Service *service.UploadService
}

func NewUploadHandler(
	service *service.UploadService,
) *UploadHandler {

	return &UploadHandler{
		Service: service,
	}
}

func (h *UploadHandler) Upload(
	c *gin.Context,
) {

	var query dto.UploadQuery

	if err := c.ShouldBindQuery(&query); err != nil {

		apperrors.Respond(
			c,
			apperrors.NewWithDetails(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid query params",
				err.Error(),
			),
		)

		return
	}

	fileHeader, err := c.FormFile("file")

	if err != nil {

		apperrors.Respond(
			c,
			apperrors.New(
				http.StatusBadRequest,
				"FILE_REQUIRED",
				"file is required",
			),
		)

		return
	}

	err = h.Service.Upload(
		c.Request.Context(),
		query.Path,
		fileHeader,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		gin.H{
			"message": "uploaded successfully",
			"file":    fileHeader.Filename,
			"size":    fileHeader.Size,
		},
	)
}
