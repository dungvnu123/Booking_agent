package auth

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/auth/register", h.Register)
	api.POST("/auth/register-patient", h.RegisterPatient)
	api.POST("/auth/login", h.Login)
	api.POST("/auth/refresh_token", h.RefreshToken)
	api.POST("/auth/logout", h.Logout)
}

func (h *Handler) Register(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Register(c.Request.Context()))
}

func (h *Handler) RegisterPatient(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.RegisterPatient(c.Request.Context()))
}

func (h *Handler) Login(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Login(c.Request.Context()))
}

func (h *Handler) RefreshToken(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.RefreshToken(c.Request.Context()))
}

func (h *Handler) Logout(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Logout(c.Request.Context()))
}
