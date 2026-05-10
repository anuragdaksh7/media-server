package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func New(
	status int,
	code string,
	message string,
) *APIError {

	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func NewWithDetails(
	status int,
	code string,
	message string,
	details interface{},
) *APIError {

	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func Respond(
	c *gin.Context,
	err error,
) {

	apiErr, ok := err.(*APIError)

	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "internal server error",
				},
			},
		)
		return
	}

	c.JSON(
		apiErr.Status,
		gin.H{
			"success": false,
			"error": gin.H{
				"code":    apiErr.Code,
				"message": apiErr.Message,
				"details": apiErr.Details,
			},
		},
	)
}
