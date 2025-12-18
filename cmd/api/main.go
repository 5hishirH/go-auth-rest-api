package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/config"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
)

func main() {
	// config setup
	// database setup
	// router setup
	// setup server

	cfg := config.MustLoad()

	userHandler := user.Handler()

	mainMux := http.NewServeMux()

	mainMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	mainMux.Handle("/api/users/", http.StripPrefix("/api/users", userHandler))

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: mainMux,
	}

	fmt.Println("Server started")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start server")
	}

}
