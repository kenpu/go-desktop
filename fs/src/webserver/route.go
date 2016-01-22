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

// this sets up the URLS
func routes() {
	r = gin.Default()

	// loads the HTML templates for rendering
	r.LoadHTMLGlob("./templates/*.html")

	// serves basic css and js files as
	// static resources
	r.StaticFS("/static", http.Dir("./static"))

	// GET /workspace/mywork renders a workspace page
	// for the workspace called "mywork"
	r.GET("/workspace/:id", nopanic(handler.Fetch))

	// GET /workspace/mywork/images/eva.jpg returns
	// the content of the file "./images/eva.jpg"
	// from the workspace directory
	r.GET("/workspace/:id/*path", nopanic(handler.Fetch))

	// WS /watch/mywork setups the websocket connection
	// which is used to inform the client that
	// some FS activity occurred.
	r.GET("/watch/:id/", nopanic(func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("[KP-debug] Watching %s", id)
		wshandler(id, c.Writer, c.Request)
	}))

	// When everything fails, render the cover page
	r.NoRoute(handler.Index)
}

// I am panicing when there is error.
// nopanic() just catches the panic and returns
// a HTTP error instead.
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
