package patient

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/patients", h.CreatePatient)
	api.GET("/patients/search", h.SearchPatients)
	api.GET("/patients/stats", h.GetPatientStats)
	api.GET("/patients", h.GetAllPatients)
	api.GET("/patients/:id", h.GetPatientById)
	api.PUT("/patients/:id", h.UpdatePatient)
	api.POST("/patients/import", h.ImportPatients)
	api.POST("/patients/import/batch", h.ImportPatientsBatch)
}

func (h *Handler) CreatePatient(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CreatePatient(c.Request.Context()))
}

func (h *Handler) SearchPatients(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.SearchPatients(c.Request.Context()))
}

func (h *Handler) GetPatientStats(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetPatientStats(c.Request.Context()))
}

func (h *Handler) GetAllPatients(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetAllPatients(c.Request.Context()))
}

func (h *Handler) GetPatientById(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetPatientById(c.Request.Context()))
}

func (h *Handler) UpdatePatient(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UpdatePatient(c.Request.Context()))
}

func (h *Handler) ImportPatients(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ImportPatients(c.Request.Context()))
}

func (h *Handler) ImportPatientsBatch(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ImportPatientsBatch(c.Request.Context()))
}
