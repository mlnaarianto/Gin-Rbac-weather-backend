package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WriteJSON adalah fungsi utama untuk mencetak format JSON secara global
func WriteJSON(c *gin.Context, code int, status string, data any) {
	c.JSON(code, gin.H{
		"code":   code,
		"status": status,
		"data":   data,
	})
}

func ResponseOK(c *gin.Context, data any) {
	WriteJSON(c, http.StatusOK, "OK", data)
}

func ResponseCreated(c *gin.Context, data any) {
	WriteJSON(c, http.StatusCreated, "Created", data)
}

func ResponseBadRequest(c *gin.Context, data any) {
	WriteJSON(c, http.StatusBadRequest, "Bad Request", data)
}

func ResponseUnauthorized(c *gin.Context, data any) {
	WriteJSON(c, http.StatusUnauthorized, "Unauthorized", data)
}

func ResponseForbidden(c *gin.Context, data any) {
	WriteJSON(c, http.StatusForbidden, "Forbidden", data)
}

func ResponseNotFound(c *gin.Context, data any) {
	WriteJSON(c, http.StatusNotFound, "Not Found", data)
}

func ResponseInternalError(c *gin.Context, data any) {
	WriteJSON(c, http.StatusInternalServerError, "Internal Server Error", data)
}