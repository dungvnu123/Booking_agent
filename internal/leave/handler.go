package leave

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/roster/leave-requests", h.Create)
	api.GET("/roster/leave-requests/me", h.ListMine)
	api.GET("/roster/leave-requests", h.ListAll)
	api.POST("/roster/leave-requests/preview", h.Preview)
	api.GET("/roster/leave-requests/:id/impact", h.Impact)
	api.POST("/roster/leave-requests/:id/approve", h.Approve)
	api.POST("/roster/leave-requests/:id/reject", h.Reject)
	api.POST("/roster/leave-requests/:id/withdraw", h.Withdraw)
}

func (h *Handler) Create(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Create(c.Request.Context()))
}

func (h *Handler) ListMine(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListMine(c.Request.Context()))
}

func (h *Handler) ListAll(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListAll(c.Request.Context()))
}

func (h *Handler) Preview(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Preview(c.Request.Context()))
}

func (h *Handler) Impact(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Impact(c.Request.Context()))
}

func (h *Handler) Approve(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Approve(c.Request.Context()))
}

func (h *Handler) Reject(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Reject(c.Request.Context()))
}

func (h *Handler) Withdraw(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Withdraw(c.Request.Context()))
}
