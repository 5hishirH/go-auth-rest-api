package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/internal/auth"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/config"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/db"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/filestore/minio"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/session"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/types"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
)

func main() {

	// config setup
	cfg := config.MustLoad()

	// filestore setup
	minioClient, err := minio.New(
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.Bucket,
		cfg.UseSSL,
	)

	if err != nil {
		log.Fatal("failed to init storage: ", err)
	}

	// database setup
	psql, err := db.NewPostgresqlStorage(cfg.DbSource)

	if err != nil {
		log.Fatal("failed to init db: ", err)
	}

	// session setup
	types.RegisterTypes()

	store, err := session.NewRedisStore(cfg)
	if err != nil {
		log.Fatal("failed to create session: ", err, "\n")
	}

	mainMux := http.NewServeMux()

	mainMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	profilePicApiPrefix := "api/user"

	// router setup
	userRepo := user.NewRepository(psql)
	authService := auth.NewService(minioClient, userRepo, "profile-pics")
	authHandler := auth.NewHandler(authService, store, cfg.Cookies.Refresh.Name, cfg.Cookies.Refresh.Path, cfg.Refresh.Expiry, cfg.Cookies.Refresh.Secure, profilePicApiPrefix)
	authRoutes := authHandler.RegisterRoutes()
	mainMux.Handle("/api/auth/", http.StripPrefix("/api/auth", authRoutes))

	authMiddleware := auth.NewMiddleware(store)
	userHandler := user.NewHandler(authMiddleware.AuthMiddleware, userRepo, profilePicApiPrefix)
	userRoutes := userHandler.RegisterRoutes()
	mainMux.Handle("/api/user/", http.StripPrefix("/api/user", userRoutes))

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: mainMux,
	}

	// setup server
	fmt.Println("Server started")
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("failed to start server")
	}

}
