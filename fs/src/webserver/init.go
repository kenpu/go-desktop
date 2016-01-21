package webserver

func Run(addr string) {
	routes()
	r.Run(addr)
}
