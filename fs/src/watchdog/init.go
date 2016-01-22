package watchdog

import (
	"log"

	"github.com/gorilla/websocket"
	"gopkg.in/fsnotify.v1"
)

// This is a struct doing some basic bookkeeping
// of the websocket connection and the physical
// filesystem directory of the workspace to be
// monitored.
type Watcher struct {
	Conn *websocket.Conn
	Dir  string
}

// a constructor for a new instance of a watcher
func New(dir string, conn *websocket.Conn) *Watcher {
	return &Watcher{
		Conn: conn,
		Dir:  dir,
	}
}

// starts the watcher.  Remember, the watcher.Conn
// talks back to the browser.
func (watcher *Watcher) Start() error {
	// fsnotify is our secret sauce here.
	// NewWatcher() starts a monitor for
	// FS activities.  It creates two channels:
	// `w.Events` and `w.Errors` for normal
	// events and errenous events.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// really cool go feature.  Now, we never keep to
	// worry about forgetting to close the websocket
	// connections.
	defer w.Close()

	// this channel just monitors the exit condition
	// of the goroutine that we are about to start.
	done := make(chan bool)

	// here we go.  We are starting the real-time
	// content authoring magic here.
	go func() {
		// a blocking loop reading the fs events
		// from the fsnotify channels.
		// I manually convert the type of events
		// to a string `method`.
		for {
			select {
			case ev := <-w.Events:
				var method string
				switch {
				case ev.Op&fsnotify.Create > 0:
					// a file is created
					method = "Create"
				case ev.Op&fsnotify.Write > 0:
					// a file is updated
					method = "Write"
				case ev.Op&fsnotify.Remove > 0:
					// a file is removed
					method = "Remove"
				case ev.Op&fsnotify.Rename > 0:
					// a file has been renamed
					// to a new name
					method = "Rename"
				case ev.Op&fsnotify.Chmod > 0:
					// the modtime is changed
					method = "Chmod"
				default:
					method = "Unknown"
				}
				// notify the client the filename and
				// the type of event through the
				// websocket
				if notify(watcher.Conn, ev.Name, method) != nil {
					break
				}
			case err = <-w.Errors:
				// when we see a fs error, we just break
				// out of the loop.
				// I am not sure if there is what we
				// should do.
				log.Printf("Modified: %s", err.Error())
				break
			}
		}

		// before exiting the goroutine, we need to
		// inform the parent that we are done.
		done <- true
	}()

	// cool magic here: tell fsnotify to start
	// monitoring the directory that is the workspace
	err = w.Add(watcher.Dir)
	if err != nil {
		log.Printf("[Watcher] %s", err.Error())
		return err
	}

	// block until the goroutine is done
	<-done

	// return.  The deferred will close the websocket
	// for us.
	return nil
}

// I am just encoding the (filename, method) as an JSON array.
func notify(conn *websocket.Conn, name, event string) error {
	data := [2]string{name, event}
	return conn.WriteJSON(data)
}
