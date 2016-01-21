package webserver

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(workspaceId string, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("wshanlder/upgrade: %+v", err)
		return
	}
	log.Printf("watching workspace %s", workspaceId)

	_ = conn
}
