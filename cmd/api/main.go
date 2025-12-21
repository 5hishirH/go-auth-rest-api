package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/config"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/filestore/minio"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
)

func main() {
	// config setup
	// storage setup
	// database setup
	// router setup
	// setup server

	cfg := config.MustLoad()

	minioClient, err := minio.NewClient(
		"localhost:9000", // Endpoint
		"minioadmin",     // Access Key
		"minioadmin",     // Secret Key
		"my-test-bucket", // Bucket Name
		false,            // SSL (false for localhost)
	)

	if err != nil {
		log.Fatal("Failed to init storage:", err)
	}

	userHandler := user.NewHandler(minioClient)
	userRoutes := userHandler.RegisterRoutes()

	mainMux := http.NewServeMux()

	mainMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	mainMux.Handle("/api/users/", http.StripPrefix("/api/users", userRoutes))

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: mainMux,
	}

	fmt.Println("Server started")
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start server")
	}

}
