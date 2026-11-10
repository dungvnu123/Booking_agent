package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

// TokenVerifier must validate HS256 signature, issuer, expiry, jti revocation
// and active user status. No permissive placeholder implementation is supplied.
type TokenVerifier interface {
	Verify(context.Context, string) (Principal, error)
}

func Authenticate(verifier TokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme, token, ok := strings.Cut(c.GetHeader("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" || verifier == nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Authentication required"})
			return
		}
		p, err := verifier.Verify(c.Request.Context(), token)
		if err != nil || p.UserID == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
			return
		}
		c.Set("principal", p)
		c.Next()
	}
}
