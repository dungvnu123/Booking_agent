package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var b [16]byte
		if _, err := rand.Read(b[:]); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		id := hex.EncodeToString(b[:])
		c.Set("requestID", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// CORS permits the API origin and explicitly configured frontend origins.
func CORS(origins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[origin] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}
		c.Header("Vary", "Origin")
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		sameOrigin := origin == scheme+"://"+c.Request.Host
		if !sameOrigin && !allowed[origin] {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type Principal struct {
	UserID string
	Roles  []string
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("principal")
		p, ok := value.(Principal)
		if !exists || !ok || strings.TrimSpace(p.UserID) == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Authentication required"})
			return
		}
		for _, role := range roles {
			for _, actual := range p.Roles {
				if role == actual {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})
	}
}
