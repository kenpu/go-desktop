package webserver

// This is invoked by the main()
func Run(addr string) {
	// this setups the routing of URL
	routes()

	// The package-local global variable
	// `r` is a "gin" engine.  It's basically
	// the web server instance.
	r.Run(addr)
}
