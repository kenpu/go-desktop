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

func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func Fetch(c *gin.Context) {
	workspaceId := c.Param("id")
	resource := c.Param("path")

	if resource == "/" || resource == "" {
		Workspace(c)
		return
	}

	log.Printf("id=\"%s\", resource=\"%s\"", workspaceId, resource)

	content, err := ioutil.ReadFile(resolve(workspaceId, resource))
	if err != nil {
		panic(err.Error())
	}

	mimetype := mime.TypeByExtension(path.Ext(resource))

	c.Data(http.StatusOK, mimetype, content)
}

func Workspace(c *gin.Context) {
	workspaceId := c.Param("id")
	c.HTML(200, "workspace.html", gin.H{
		"Id": workspaceId,
	})
}

func resolve(id, p string) string {
	return filepath.Join("./workspaces", id, p)
}
