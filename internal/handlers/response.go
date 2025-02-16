package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"errors"`
}

func newErrorResponse(context *gin.Context, statusCode int, message string) {
	slog.Error(message)
	context.AbortWithStatusJSON(statusCode, Error{message})
}
