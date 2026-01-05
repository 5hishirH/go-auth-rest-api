package auth

import "net/http"

func (h Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", h.Register)
	return mux
}
