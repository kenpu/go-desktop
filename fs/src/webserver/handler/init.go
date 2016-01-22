package handler

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// frontpage when no routes are possible.
func Index(c *gin.Context) {
	// this is how gin API renders
	// response.  c.HTML will render
	// go template for us.
	// gin.H is a shorthand for map[string]interface{}
	c.HTML(200, "index.html", gin.H{})
}

// this is when we return raw data back to the client
// for fetching files from a workspace
func Fetch(c *gin.Context) {
	// this is the workspace ID
	workspaceId := c.Param("id")
	// this is the relative path of the resource
	resource := c.Param("path")

	// if the resource is empty, then
	// we will render a workspace response
	// see `Workspace`
	if resource == "/" || resource == "" {
		Workspace(c)
		return
	}

	log.Printf("id=\"%s\", resource=\"%s\"", workspaceId, resource)

	// slurp up the entire file, and return it
	// as mime content-type data
	content, err := ioutil.ReadFile(Resolve(workspaceId, resource))
	if err != nil {
		panic(err.Error())
	}

	// use `mime.TypeByExtension` to guess the content-type
	mimetype := mime.TypeByExtension(path.Ext(resource))

	// again gin API is really nice and simple
	c.Data(http.StatusOK, mimetype, content)
}

// this is called with the resource is empty.
// We render a template.
func Workspace(c *gin.Context) {
	workspaceId := c.Param("id")
	c.HTML(200, "workspace.html", gin.H{
		"Id": workspaceId,
	})
}

// this helps to resolve a workspace ID
// and a relative path to a physical filesystem
// directory.
// TODO: we should read in some config variable
// of where workspace root should be.
func Resolve(id, p string) string {
	return filepath.Join("./workspaces", id, p)
}
