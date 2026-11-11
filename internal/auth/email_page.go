package auth

import (
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"html/template"
	"net/http"

	"abclinic.local/backend/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

//go:embed verification.html
var verificationHTML string

var verificationPage = template.Must(template.New("email-verification").Parse(verificationHTML))

func (h *Handler) EmailVerificationPage(c *gin.Context) {
	var bytes [24]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		httpx.Fail(c, err)
		return
	}
	nonce := base64.RawStdEncoding.EncodeToString(bytes[:])
	c.Header("Cache-Control", "no-store")
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'none'; script-src 'nonce-"+nonce+"'; style-src 'nonce-"+nonce+"'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'none'")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	if err := verificationPage.Execute(c.Writer, struct{ Nonce string }{nonce}); err != nil {
		_ = c.Error(err)
	}
}
