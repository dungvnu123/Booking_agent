package swagger

import (
	"abclinic.local/backend/api"
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:embed assets/*
var files embed.FS

func Register(r *gin.Engine) {
	redirect := func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html") }
	r.GET("/", redirect)
	r.GET("/swagger", redirect)
	r.GET("/api/swagger-ui.html", redirect)
	r.GET("/api/v3/api-docs", func(c *gin.Context) { c.Data(200, "application/json; charset=utf-8", api.Document) })
	r.GET("/swagger/*file", func(c *gin.Context) {
		file := c.Param("file")
		var contentType string
		switch file {
		case "/", "/index.html":
			file, contentType = "/index.html", "text/html; charset=utf-8"
		case "/openapi.json":
			c.Data(200, "application/json; charset=utf-8", api.Document)
			return
		case "/swagger-ui.css":
			contentType = "text/css; charset=utf-8"
		case "/swagger-ui-bundle.js", "/swagger-initializer.js":
			contentType = "application/javascript; charset=utf-8"
		default:
			c.JSON(404, gin.H{"message": "Not found"})
			return
		}
		body, err := files.ReadFile("assets" + file)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(200, contentType, body)
	})
}
