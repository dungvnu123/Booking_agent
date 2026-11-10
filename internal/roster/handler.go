package roster

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/roster/monthly/:id/picks/me", h.GetMyPick)
	api.PUT("/roster/monthly/:id/picks/me", h.SubmitPick)
	api.POST("/roster/monthly", h.OpenMonth)
	api.GET("/roster/monthly", h.ListMonths)
	api.GET("/roster/monthly/current", h.GetCurrentOpenMonth)
	api.GET("/roster/monthly/my-preview", h.GetMyPreview)
	api.GET("/roster/monthly/:id", h.GetMonth)
	api.GET("/roster/monthly/capacity-config", h.GetAvgCapacity)
	api.PUT("/roster/monthly/solver-config", h.UpdateSolverConfig)
	api.PUT("/roster/monthly/:id/override-config", h.OverrideConfig)
	api.DELETE("/roster/monthly/:id/override-config", h.ClearOverrideConfig)
	api.POST("/roster/monthly/:id/solve", h.RunSolver)
	api.GET("/roster/monthly/:id/result", h.GetResult)
	api.GET("/roster/monthly/:id/publish-preview", h.GetPublishPreview)
	api.POST("/roster/monthly/:id/approve", h.Approve)
	api.POST("/roster/monthly/:id/revert-to-draft", h.RevertToDraft)
	api.POST("/roster/monthly/:id/publish", h.Publish)
	api.POST("/roster/monthly/:id/publish-now", h.PublishNow)
	api.POST("/roster/monthly/:id/assignments", h.AddAssignment)
	api.DELETE("/roster/monthly/:id/assignments", h.RemoveAssignment)
	api.POST("/roster/monthly/:id/assignments/lock", h.LockAssignment)
	api.POST("/roster/monthly/:id/lock", h.LockCell)
	api.GET("/roster/monthly/:id/cell-candidates", h.ListCellCandidates)
	api.GET("/roster/monthly/:id/doctor-statuses", h.ListDoctorPickStatuses)
	api.POST("/roster/monthly/:id/doctors/:doctorId/reopen", h.ReopenForDoctor)
	api.POST("/roster/monthly/:id/doctors/:doctorId/lock", h.LockForDoctor)
	api.POST("/roster/monthly/:id/doctors/:doctorId/unlock", h.UnlockForDoctor)
	api.GET("/roster/monthly/:id/published-swap/impact", h.GetPublishedSwapImpact)
	api.POST("/roster/monthly/:id/published-swap", h.SwapPublishedDoctor)
	api.GET("/roster/monthly/:id/audit-log", h.ListAuditLog)
	api.GET("/roster/work-schedules", h.ListWorkSchedules)
	api.GET("/roster/work-schedules/me", h.ListMyWorkSchedules)
}

func (h *Handler) GetMyPick(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetMyPick(c.Request.Context()))
}

func (h *Handler) SubmitPick(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.SubmitPick(c.Request.Context()))
}

func (h *Handler) OpenMonth(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.OpenMonth(c.Request.Context()))
}

func (h *Handler) ListMonths(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListMonths(c.Request.Context()))
}

func (h *Handler) GetCurrentOpenMonth(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetCurrentOpenMonth(c.Request.Context()))
}

func (h *Handler) GetMyPreview(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetMyPreview(c.Request.Context()))
}

func (h *Handler) GetMonth(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetMonth(c.Request.Context()))
}

func (h *Handler) GetAvgCapacity(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetAvgCapacity(c.Request.Context()))
}

func (h *Handler) UpdateSolverConfig(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UpdateSolverConfig(c.Request.Context()))
}

func (h *Handler) OverrideConfig(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.OverrideConfig(c.Request.Context()))
}

func (h *Handler) ClearOverrideConfig(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ClearOverrideConfig(c.Request.Context()))
}

func (h *Handler) RunSolver(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.RunSolver(c.Request.Context()))
}

func (h *Handler) GetResult(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetResult(c.Request.Context()))
}

func (h *Handler) GetPublishPreview(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetPublishPreview(c.Request.Context()))
}

func (h *Handler) Approve(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Approve(c.Request.Context()))
}

func (h *Handler) RevertToDraft(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.RevertToDraft(c.Request.Context()))
}

func (h *Handler) Publish(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Publish(c.Request.Context()))
}

func (h *Handler) PublishNow(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.PublishNow(c.Request.Context()))
}

func (h *Handler) AddAssignment(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.AddAssignment(c.Request.Context()))
}

func (h *Handler) RemoveAssignment(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.RemoveAssignment(c.Request.Context()))
}

func (h *Handler) LockAssignment(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.LockAssignment(c.Request.Context()))
}

func (h *Handler) LockCell(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.LockCell(c.Request.Context()))
}

func (h *Handler) ListCellCandidates(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListCellCandidates(c.Request.Context()))
}

func (h *Handler) ListDoctorPickStatuses(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListDoctorPickStatuses(c.Request.Context()))
}

func (h *Handler) ReopenForDoctor(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ReopenForDoctor(c.Request.Context()))
}

func (h *Handler) LockForDoctor(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.LockForDoctor(c.Request.Context()))
}

func (h *Handler) UnlockForDoctor(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UnlockForDoctor(c.Request.Context()))
}

func (h *Handler) GetPublishedSwapImpact(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetPublishedSwapImpact(c.Request.Context()))
}

func (h *Handler) SwapPublishedDoctor(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.SwapPublishedDoctor(c.Request.Context()))
}

func (h *Handler) ListAuditLog(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListAuditLog(c.Request.Context()))
}

func (h *Handler) ListWorkSchedules(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListWorkSchedules(c.Request.Context()))
}

func (h *Handler) ListMyWorkSchedules(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListMyWorkSchedules(c.Request.Context()))
}
