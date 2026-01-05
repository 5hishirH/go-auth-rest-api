package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/config"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisStore(cfg *config.Config) (*Store, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	maxAgeDuration, err := time.ParseDuration(cfg.Cookies.Session.Expiry)
	if err != nil {
		return nil, fmt.Errorf("invalid session expiry format: %w", err)
	}

	store.Options(sessions.Options{
		Path:     cfg.Cookies.Session.Path,
		MaxAge:   int(maxAgeDuration.Seconds()),
		HttpOnly: true,
		Secure:   cfg.Cookies.Session.Secure,
		SameSite: http.SameSiteNoneMode,
	})

	store.KeyPrefix("session:")

	return &Store{
		inner:      store,
		cookieName: cfg.Cookies.Session.Name,
	}, nil
}
