package router

import "net/http"

// Router representation of a router
type Router interface {
	POST(uri string, f func(w http.ResponseWriter, r *http.Request))
	GET(uri string, f func(w http.ResponseWriter, r *http.Request))
	SERVE(port string)
}
