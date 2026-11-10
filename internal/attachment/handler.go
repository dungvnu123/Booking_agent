package attachment

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/files/upload", h.UploadFile)
	api.GET("/files/list", h.ListFiles)
}

func (h *Handler) UploadFile(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UploadFile(c.Request.Context()))
}

func (h *Handler) ListFiles(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListFiles(c.Request.Context()))
}
