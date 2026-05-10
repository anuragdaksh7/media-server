package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	apperrors "fileserver/pkg/errors"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {

		defer func() {

			if err := recover(); err != nil {

				requestID := c.GetString(RequestIDKey)

				slog.Error(
					"panic recovered",
					slog.String(
						"request_id",
						requestID,
					),
					slog.String(
						"panic",
						fmt.Sprintf("%v", err),
					),
					slog.String(
						"stack_trace",
						string(debug.Stack()),
					),
				)

				apperrors.Respond(
					c,
					apperrors.New(
						http.StatusInternalServerError,
						"INTERNAL_SERVER_ERROR",
						"internal server error",
					),
				)

				c.Abort()
			}
		}()

		c.Next()
	}
}
