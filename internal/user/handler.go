package user

import (
	"log/slog"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("creating a user")

	w.Write([]byte("Welcome to user api"))
}
