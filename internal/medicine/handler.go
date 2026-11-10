package medicine

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/medicines", h.GetMedicines)
	api.GET("/medicines/metadata", h.GetMetadata)
	api.GET("/medicines/:id", h.GetMedicineById)
	api.POST("/medicines", h.CreateMedicine)
	api.PUT("/medicines/:id", h.UpdateMedicine)
	api.DELETE("/medicines/:id", h.DeleteMedicine)
	api.POST("/medicines/bulk", h.BulkImport)
}

func (h *Handler) GetMedicines(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetMedicines(c.Request.Context()))
}

func (h *Handler) GetMetadata(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetMetadata(c.Request.Context()))
}

func (h *Handler) GetMedicineById(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetMedicineById(c.Request.Context()))
}

func (h *Handler) CreateMedicine(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CreateMedicine(c.Request.Context()))
}

func (h *Handler) UpdateMedicine(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UpdateMedicine(c.Request.Context()))
}

func (h *Handler) DeleteMedicine(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.DeleteMedicine(c.Request.Context()))
}

func (h *Handler) BulkImport(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.BulkImport(c.Request.Context()))
}
