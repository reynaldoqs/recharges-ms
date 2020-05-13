package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type chiRouter struct{}

var chiDispatcher = chi.NewRouter()

// NewChiRouter creates a new chi router
func NewChiRouter() Router {
	chiDispatcher.Use(middleware.RequestID)
	chiDispatcher.Use(middleware.RealIP)
	chiDispatcher.Use(middleware.Logger)
	chiDispatcher.Use(middleware.Recoverer)
	return &chiRouter{}
}

func (*chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, f)
}

func (*chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, f)
}

func (*chiRouter) SERVE(port string) {
	fmt.Printf("Chi HTTP server running on port %v\n", port)
	http.ListenAndServe(port, chiDispatcher)
}
