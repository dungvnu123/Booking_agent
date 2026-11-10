package examination

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/examinations/:patientId", h.Create)
	api.GET("/examinations/patient/:patientId", h.GetByPatient)
	api.GET("/examinations/:id", h.GetById)
	api.PUT("/examinations/:id", h.Update)
}

func (h *Handler) Create(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Create(c.Request.Context()))
}

func (h *Handler) GetByPatient(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetByPatient(c.Request.Context()))
}

func (h *Handler) GetById(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetById(c.Request.Context()))
}

func (h *Handler) Update(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Update(c.Request.Context()))
}
