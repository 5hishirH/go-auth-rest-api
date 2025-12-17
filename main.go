package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/config"
)

func main() {
	// config setup
	// database setup
	// router setup
	// setup server

	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Println("Server started")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start server")
	}

}
