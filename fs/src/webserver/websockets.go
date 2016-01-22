package webserver

import (
	"log"
	"net/http"
	"watchdog"
	"webserver/handler"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// this is called by the routing setup on /watch/<workspace>
// it resolves the workspace ID to its physical directory
// and then starts a `Watchdog` which is actually just
// a goroutine that monitors the file system events.
func wshandler(workspaceId string, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("wshanlder/upgrade: %+v", err)
		return
	}
	log.Printf("watching workspace %s", workspaceId)

	dir := handler.Resolve(workspaceId, "")

	// this is where the watchdog starts.
	// nothing is closed here yet.
	watchdog.New(dir, conn).Start()
}
