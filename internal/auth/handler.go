package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

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
	api.GET("/auth/verify-email", h.EmailVerificationPage)
	api.POST("/auth/verify-email", h.VerifyEmail)
	api.POST("/auth/resend-verification", h.ResendVerification)
}

func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequest
	if !bind(c, &request) {
		return
	}
	user, err := h.service.Register(c.Request.Context(), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusCreated, "Đã tạo tài khoản. Vui lòng kiểm tra email để xác thực.", user)
}

func (h *Handler) RegisterPatient(c *gin.Context) {
	var request RegisterPatientRequest
	if !bind(c, &request) {
		return
	}
	user, err := h.service.RegisterPatient(c.Request.Context(), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusCreated, "Đã tạo tài khoản bệnh nhân. Vui lòng kiểm tra email để xác thực.", user)
}

func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest
	if !bind(c, &request) {
		return
	}
	result, err := h.service.Login(c.Request.Context(), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusOK, "Login successful!", result)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var request TokenRequest
	if !bind(c, &request) {
		return
	}
	result, err := h.service.RefreshToken(c.Request.Context(), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusOK, "Refresh token successful!", result)
}

func (h *Handler) Logout(c *gin.Context) {
	var request TokenRequest
	if !bind(c, &request) {
		return
	}
	if err := h.service.Logout(c.Request.Context(), request); err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusOK, "Logout successful!", nil)
}

func bind(c *gin.Context, destination any) bool {
	c.Header("Cache-Control", "no-store")
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16<<10)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		httpx.Fail(c, ErrInvalidInput)
		return false
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		httpx.Fail(c, ErrInvalidInput)
		return false
	}
	return true
}

func respond(c *gin.Context, status int, message string, data any) {
	c.Header("Cache-Control", "no-store")
	c.JSON(status, httpx.Response{StatusCode: 1000, Message: message, Data: data})
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var request VerifyEmailRequest
	if !bind(c, &request) {
		return
	}
	if err := h.service.VerifyEmail(c.Request.Context(), request); err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusOK, "Email đã được xác thực. Bạn có thể đăng nhập.", nil)
}

func (h *Handler) ResendVerification(c *gin.Context) {
	var request ResendVerificationRequest
	if !bind(c, &request) {
		return
	}
	if err := h.service.ResendVerification(c.Request.Context(), request); err != nil {
		httpx.Fail(c, err)
		return
	}
	respond(c, http.StatusAccepted, "Nếu tài khoản cần xác thực và đủ thời gian gửi lại, email xác thực sẽ được gửi.", nil)
}
