package notification

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/notifications/me", h.ListMine)
	api.POST("/notifications/:id/read", h.MarkRead)
}

func (h *Handler) ListMine(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListMine(c.Request.Context()))
}

func (h *Handler) MarkRead(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.MarkRead(c.Request.Context()))
}
