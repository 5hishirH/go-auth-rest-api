package user

import "net/http"

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /profile", h.requireAuth(h.Profile))
	return mux
}
