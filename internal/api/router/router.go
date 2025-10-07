package router

import (
	"github.com/K1la/image-processor/internal/api/handler"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wb-go/wbf/ginext"
)

func New(handler *handler.Handler) *ginext.Engine {
	e := ginext.New()
	e.Use(ginext.Recovery(), ginext.Logger())

	// API routes first
	api := e.Group("/api/")
	{
		api.POST("/upload", handler.CreateImage)
		api.GET("/image/:id", handler.GetImageByID)
		api.GET("/image/info/:id", handler.GetImageInfoByID)
		api.DELETE("/image/:id", handler.DeleteImage)
	}

	// Frontend: serve files from ./web without conflicting wildcard
	e.NoRoute(func(c *ginext.Context) {
		if c.Request.URL.Path == "/" {
			http.ServeFile(c.Writer, c.Request, "./web/index.html")
			return
		}
		safe := filepath.Clean("." + c.Request.URL.Path)
		filePath := filepath.Join("./web", safe)
		if fi, err := os.Stat(filePath); err == nil && !fi.IsDir() {
			http.ServeFile(c.Writer, c.Request, filePath)
			return
		}
		c.Status(http.StatusNotFound)
	})

	return e
}
