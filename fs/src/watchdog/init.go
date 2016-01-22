package watchdog

import (
	"log"

	"github.com/gorilla/websocket"
	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	Conn *websocket.Conn
	Dir  string
}

func New(dir string, conn *websocket.Conn) *Watcher {
	return &Watcher{
		Conn: conn,
		Dir:  dir,
	}
}

func (watcher *Watcher) Start() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-w.Events:
				var method string
				switch {
				case ev.Op&fsnotify.Create > 0:
					method = "Create"
				case ev.Op&fsnotify.Write > 0:
					method = "Write"
				case ev.Op&fsnotify.Remove > 0:
					method = "Remove"
				case ev.Op&fsnotify.Rename > 0:
					method = "Rename"
				case ev.Op&fsnotify.Chmod > 0:
					method = "Chmod"
				default:
					method = "Unknown"
				}
				log.Printf("%s: %s", method, ev.Name)
				if notify(watcher.Conn, ev.Name, method) != nil {
					break
				}
			case err = <-w.Errors:
				log.Printf("Modified: %s", err.Error())
				break
			}
		}
		done <- true
	}()

	err = w.Add(watcher.Dir)
	if err != nil {
		log.Printf("[Watcher] %s", err.Error())
		return err
	}

	<-done

	return nil
}

func notify(conn *websocket.Conn, name, event string) error {
	data := [2]string{name, event}
	return conn.WriteJSON(data)
}
