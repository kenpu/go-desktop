package webserver

import (
	"log"
	"net/http"
	"webserver/handler"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
)

func routes() {
	r = gin.Default()

	r.LoadHTMLGlob("./templates/*.html")

	r.StaticFS("/static", http.Dir("./static"))

	r.GET("/workspace/:id", nopanic(handler.Fetch))
	r.GET("/workspace/:id/*path", nopanic(handler.Fetch))

	r.GET("/watch/:id/", nopanic(func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("[KP-debug] Watching %s", id)
		wshandler(id, c.Writer, c.Request)
	}))

	r.NoRoute(handler.Index)
}

func nopanic(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.String(http.StatusNotFound, "Err: %s", r)
			}
		}()

		f(c)
	}
}
