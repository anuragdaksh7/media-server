package handler

import (
	"net/http"

	"fileserver/internal/filesystem/dto"
	"fileserver/internal/filesystem/service"

	apperrors "fileserver/pkg/errors"
	"fileserver/pkg/response"

	"github.com/gin-gonic/gin"
)

type FilesystemHandler struct {
	Service *service.FilesystemService
}

func NewFilesystemHandler(
	service *service.FilesystemService,
) *FilesystemHandler {

	return &FilesystemHandler{
		Service: service,
	}
}

func (h *FilesystemHandler) List(
	c *gin.Context,
) {

	var query dto.ListQuery

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

	files, err := h.Service.List(
		c.Request.Context(),
		query.Path,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		files,
	)
}

func (h *FilesystemHandler) Stat(
	c *gin.Context,
) {

	var query dto.StatQuery

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

	file, err := h.Service.Stat(
		c.Request.Context(),
		query.Path,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		file,
	)
}

func (h *FilesystemHandler) Delete(
	c *gin.Context,
) {

	var req dto.DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		apperrors.Respond(
			c,
			apperrors.NewWithDetails(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid request body",
				err.Error(),
			),
		)

		return
	}

	err := h.Service.Delete(
		c.Request.Context(),
		req.Path,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "deleted successfully",
		},
	)
}

func (h *FilesystemHandler) Move(
	c *gin.Context,
) {

	var req dto.MoveRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		apperrors.Respond(
			c,
			apperrors.NewWithDetails(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid request body",
				err.Error(),
			),
		)

		return
	}

	err := h.Service.Move(
		c.Request.Context(),
		req.Source,
		req.Destination,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "moved successfully",
		},
	)
}

func (h *FilesystemHandler) Copy(
	c *gin.Context,
) {

	var req dto.CopyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		apperrors.Respond(
			c,
			apperrors.NewWithDetails(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid request body",
				err.Error(),
			),
		)

		return
	}

	err := h.Service.Copy(
		c.Request.Context(),
		req.Source,
		req.Destination,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "copied successfully",
		},
	)
}

func (h *FilesystemHandler) Mkdir(
	c *gin.Context,
) {

	var req dto.MkdirRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		apperrors.Respond(
			c,
			apperrors.NewWithDetails(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid request body",
				err.Error(),
			),
		)

		return
	}

	err := h.Service.Mkdir(
		c.Request.Context(),
		req.Path,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		gin.H{
			"message": "directory created",
		},
	)
}
