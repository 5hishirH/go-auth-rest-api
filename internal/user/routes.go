package user

import (
	"net/http"
)

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", Create)

	return mux
}
