package handler

import (
	"fileserver/internal/streaming/dto"
	"fileserver/internal/streaming/service"
	"net/http"

	apperrors "fileserver/pkg/errors"

	"github.com/gin-gonic/gin"
)

type StreamingHandler struct {
	Service *service.StreamingService
}

func NewStreamingHandler(
	service *service.StreamingService,
) *StreamingHandler {

	return &StreamingHandler{
		Service: service,
	}
}

func (h *StreamingHandler) Stream(
	c *gin.Context,
) {

	var query dto.StreamQuery

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

	file, info, err := h.Service.OpenFile(
		c.Request.Context(),
		query.Path,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	defer file.Close()

	c.Header(
		"Content-Type",
		info.MimeType,
	)

	c.Header(
		"Accept-Ranges",
		"bytes",
	)

	c.Header(
		"Cache-Control",
		"public, max-age=3600",
	)

	http.ServeContent(
		c.Writer,
		c.Request,
		info.Name,
		info.ModifiedAt,
		file,
	)
}

func (h *StreamingHandler) Download(
	c *gin.Context,
) {

	var query dto.DownloadQuery

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

	file, info, err := h.Service.OpenFile(
		c.Request.Context(),
		query.Path,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	defer file.Close()

	c.Header(
		"Content-Disposition",
		`attachment; filename="`+info.Name+`"`,
	)

	c.Header(
		"Content-Type",
		"application/octet-stream",
	)

	c.Header(
		"Content-Length",
		string(rune(info.Size)),
	)

	http.ServeContent(
		c.Writer,
		c.Request,
		info.Name,
		info.ModifiedAt,
		file,
	)
}
