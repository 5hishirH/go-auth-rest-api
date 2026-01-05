package user

// import (
// 	"net/http"

// 	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/filestore"
// )

// type Handler struct {
// 	fileStore filestore.FileStore
// 	// service *Service (Add your service here later)
// }

// func NewHandler(fs filestore.FileStore) *Handler {
// 	return &Handler{fileStore: fs}
// }

// // RegisterRoutes registers routes to a mux
// func (h *Handler) RegisterRoutes() http.Handler {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("POST /{$}", h.Create) // Notice: h.Create
// 	return mux
// }
