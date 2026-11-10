package user

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/admin/users/roles", h.GetRoles)
	api.GET("/admin/users", h.GetUsers)
	api.GET("/admin/users/search", h.SearchUsers)
	api.GET("/admin/users/:id", h.GetUser)
	api.POST("/admin/users", h.CreateUser)
	api.PUT("/admin/users/:id", h.UpdateUser)
	api.PATCH("/admin/users/:id/status", h.ChangeStatus)
	api.PATCH("/admin/users/:id/password", h.ResetPassword)
}

func (h *Handler) GetRoles(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetRoles(c.Request.Context()))
}

func (h *Handler) GetUsers(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetUsers(c.Request.Context()))
}

func (h *Handler) SearchUsers(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.SearchUsers(c.Request.Context()))
}

func (h *Handler) GetUser(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetUser(c.Request.Context()))
}

func (h *Handler) CreateUser(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CreateUser(c.Request.Context()))
}

func (h *Handler) UpdateUser(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UpdateUser(c.Request.Context()))
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ChangeStatus(c.Request.Context()))
}

func (h *Handler) ResetPassword(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ResetPassword(c.Request.Context()))
}
