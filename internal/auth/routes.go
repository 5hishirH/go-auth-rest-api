package auth

import "net/http"

func (h Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /logout", h.Logout)
	return mux
}
