package server

import (
	"goproxy/internal/application/controllers"
	"net/http"
)

func NewServeMux(handler *controllers.ProxyController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	return mux
}
