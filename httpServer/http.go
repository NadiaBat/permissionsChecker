package httpServer

import (
	"net/http"
	"sync"
)

type Http struct { }
type httpHandler struct {wg *sync.WaitGroup}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	switch r.URL.Path {
	case "/check/":
		println("!!!!")
	default:
		println(r.URL.Path)
	}
}

func (h Http) ServeHttp(wg *sync.WaitGroup)  {
	handler := httpHandler{wg: wg}
	http.ListenAndServe(":9999", handler)
}
