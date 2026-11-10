package closure

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/clinic-closures", h.ListByMonth)
	api.GET("/clinic-closures/impact", h.GetImpact)
	api.GET("/clinic-closures/impact-bulk", h.GetImpactBulk)
	api.POST("/clinic-closures", h.Add)
	api.POST("/clinic-closures/bulk", h.AddBulk)
	api.DELETE("/clinic-closures/:id", h.Delete)
}

func (h *Handler) ListByMonth(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListByMonth(c.Request.Context()))
}

func (h *Handler) GetImpact(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetImpact(c.Request.Context()))
}

func (h *Handler) GetImpactBulk(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetImpactBulk(c.Request.Context()))
}

func (h *Handler) Add(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Add(c.Request.Context()))
}

func (h *Handler) AddBulk(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.AddBulk(c.Request.Context()))
}

func (h *Handler) Delete(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Delete(c.Request.Context()))
}
