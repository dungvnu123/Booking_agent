package booking

import (
	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/booking/available-slots", h.GetAvailableSlots)
	api.GET("/booking/slot-board", h.GetSlotBoard)
	api.POST("/booking/hold", h.HoldSlot)
	api.POST("/booking/release-hold", h.ReleaseHold)
	api.POST("/booking", h.CreateBooking)
	api.POST("/booking/staff", h.CreateStaffBooking)
	api.GET("/booking/me", h.ListMyBookings)
	api.POST("/booking/:id/cancel", h.CancelBooking)
	api.GET("/booking/work-schedules/:id", h.ListBookingsForWorkSchedule)
	api.POST("/booking/:id/check-in", h.CheckIn)
	api.POST("/booking/:id/undo-check-in", h.UndoCheckIn)
	api.POST("/booking/:id/no-show", h.MarkNoShow)
	api.POST("/booking/:id/call", h.CallBooking)
	api.POST("/booking/:id/finish", h.FinishBooking)
	api.POST("/booking/:id/call-failed", h.CallBookingFailed)
	api.POST("/booking/:id/staff-cancel", h.StaffCancelBooking)
	api.GET("/booking/:id/alternative-doctors", h.GetAlternativeDoctors)
	api.GET("/booking/:id/staff-alternative-doctors", h.GetAlternativeDoctorsForStaff)
	api.POST("/booking/:id/switch-doctor", h.SwitchDoctor)
	api.POST("/booking/:id/staff-switch-doctor", h.StaffSwitchDoctor)
	api.POST("/booking/:id/reschedule", h.Reschedule)
	api.POST("/booking/:id/staff-reschedule", h.StaffReschedule)
	api.GET("/booking/queue", h.GetQueue)
	api.POST("/booking/:id/swap", h.ManualSwap)
}

func (h *Handler) GetAvailableSlots(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetAvailableSlots(c.Request.Context()))
}

func (h *Handler) GetSlotBoard(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetSlotBoard(c.Request.Context()))
}

func (h *Handler) HoldSlot(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.HoldSlot(c.Request.Context()))
}

func (h *Handler) ReleaseHold(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ReleaseHold(c.Request.Context()))
}

func (h *Handler) CreateBooking(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CreateBooking(c.Request.Context()))
}

func (h *Handler) CreateStaffBooking(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CreateStaffBooking(c.Request.Context()))
}

func (h *Handler) ListMyBookings(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListMyBookings(c.Request.Context()))
}

func (h *Handler) CancelBooking(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CancelBooking(c.Request.Context()))
}

func (h *Handler) ListBookingsForWorkSchedule(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ListBookingsForWorkSchedule(c.Request.Context()))
}

func (h *Handler) CheckIn(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CheckIn(c.Request.Context()))
}

func (h *Handler) UndoCheckIn(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.UndoCheckIn(c.Request.Context()))
}

func (h *Handler) MarkNoShow(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.MarkNoShow(c.Request.Context()))
}

func (h *Handler) CallBooking(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CallBooking(c.Request.Context()))
}

func (h *Handler) FinishBooking(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.FinishBooking(c.Request.Context()))
}

func (h *Handler) CallBookingFailed(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.CallBookingFailed(c.Request.Context()))
}

func (h *Handler) StaffCancelBooking(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.StaffCancelBooking(c.Request.Context()))
}

func (h *Handler) GetAlternativeDoctors(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetAlternativeDoctors(c.Request.Context()))
}

func (h *Handler) GetAlternativeDoctorsForStaff(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetAlternativeDoctorsForStaff(c.Request.Context()))
}

func (h *Handler) SwitchDoctor(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.SwitchDoctor(c.Request.Context()))
}

func (h *Handler) StaffSwitchDoctor(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.StaffSwitchDoctor(c.Request.Context()))
}

func (h *Handler) Reschedule(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.Reschedule(c.Request.Context()))
}

func (h *Handler) StaffReschedule(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.StaffReschedule(c.Request.Context()))
}

func (h *Handler) GetQueue(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.GetQueue(c.Request.Context()))
}

func (h *Handler) ManualSwap(c *gin.Context) {
	// Bind a dedicated request DTO and enforce identity/ownership when implementing.
	httpx.Fail(c, h.service.ManualSwap(c.Request.Context()))
}
