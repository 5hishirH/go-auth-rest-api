package session

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/config"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
)

func NewPostgresStore(db *sql.DB, cfg *config.Config) (*Store, error) {

	key := []byte(cfg.Cookies.Session.SecretKey)

	// pgstore saves sessions to the 'http_sessions' table in Postgres
	pgStore, err := pgstore.NewPGStoreFromPool(db, key)
	if err != nil {
		return nil, err
	}

	maxAgeDuration, err := time.ParseDuration(cfg.Cookies.Session.Expiry)
	if err != nil {
		return nil, fmt.Errorf("invalid session expiry format: %w", err)
	}

	pgStore.Options = &sessions.Options{
		Path:     cfg.Cookies.Session.Path,
		MaxAge:   int(maxAgeDuration),
		HttpOnly: true,
		Secure:   cfg.Cookies.Session.Secure,
		SameSite: http.SameSiteLaxMode,
	}

	return &Store{inner: pgStore, cookieName: cfg.Cookies.Session.Name}, nil
}
