package handler

import (
	"fileserver/internal/torrent/dto"
	"fileserver/internal/torrent/service"
	"net/http"

	apperrors "fileserver/pkg/errors"
	"fileserver/pkg/response"

	"github.com/gin-gonic/gin"
)

type TorrentHandler struct {
	Service *service.TorrentService
}

func NewTorrentHandler(
	service *service.TorrentService,
) *TorrentHandler {

	return &TorrentHandler{
		Service: service,
	}
}

func (h *TorrentHandler) AddTorrent(
	c *gin.Context,
) {

	var req dto.AddTorrentRequest

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

	torrent, err := h.Service.AddTorrent(
		c.Request.Context(),
		req,
	)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		torrent,
	)
}
