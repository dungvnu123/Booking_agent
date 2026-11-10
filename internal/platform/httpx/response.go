package httpx

import (
	"abclinic.local/backend/internal/platform/fault"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response is the common wrapper. Examination and attachment contracts must be
// checked separately against the Java source before implementing their output.
type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) { c.JSON(http.StatusOK, Response{1000, "OK", data}) }
func Fail(c *gin.Context, err error) {
	if errors.Is(err, fault.ErrNotImplemented) {
		c.JSON(http.StatusNotImplemented, Response{1501, err.Error(), nil})
		return
	}
	c.JSON(http.StatusInternalServerError, Response{1500, "Internal server error", nil})
}

// Scaffold codes 1500/1501 are provisional, not verified Java business codes.
