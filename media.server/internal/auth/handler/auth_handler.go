package handler

import (
	"fileserver/internal/auth/dto"
	"fileserver/internal/auth/service"
	apperrors "fileserver/pkg/errors"
	"fileserver/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(
	service *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

func (h *AuthHandler) Register(
	c *gin.Context,
) {

	var req dto.RegisterRequest

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

	token, err := h.Service.Register(req)

	if err != nil {
		apperrors.Respond(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		dto.AuthResponse{
			Token: token,
		},
	)
}

func (h *AuthHandler) Login(
	c *gin.Context,
) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	token, err := h.Service.Login(req)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.AuthResponse{
			Token: token,
		},
	)
}
