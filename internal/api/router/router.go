package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/K1la/image-processor/internal/api/handler"

	"github.com/wb-go/wbf/ginext"
)

func New(handler *handler.Handler) *ginext.Engine {
	e := ginext.New()
	e.Use(ginext.Recovery(), ginext.Logger())

	// API routes first
	api := e.Group("/api/")
	{
		api.POST("/upload", handler.CreateImage)
		api.GET("/image/:id", handler.GetImageById)
		api.GET("/image/info/:id", handler.GetImageInfoByID)
		api.DELETE("/image/:id", handler.DeleteImageByID)

	}

	// Frontend: serve files from ./web without conflicting wildcard
	e.NoRoute(func(c *ginext.Context) {
		if c.Request.URL.Path == "/" {
			http.ServeFile(c.Writer, c.Request, "./web/index.html")
			return
		}
		// Serve only files under /web/ directly from disk
		if strings.HasPrefix(c.Request.URL.Path, "/web/") {
			safe := filepath.Clean("." + c.Request.URL.Path)
			http.ServeFile(c.Writer, c.Request, safe)
			return
		}
		c.Status(http.StatusNotFound)
	})

	return e
}
