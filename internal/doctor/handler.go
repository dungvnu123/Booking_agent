package doctor

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/admin/doctor-profiles", h.ListDoctorProfiles)
	api.POST("/admin/doctor-profiles", h.AssignProfile)
}

func (h *Handler) ListDoctorProfiles(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListDoctorProfiles(c.Request.Context()))
}

func (h *Handler) AssignProfile(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.AssignProfile(c.Request.Context()))
}
