package auth

import (
	"context"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/response"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/session"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/types"
)

type Middleware struct {
	sessionStore *session.Store
}

func NewMiddleware(sessionStore *session.Store) *Middleware {
	return &Middleware{
		sessionStore: sessionStore,
	}
}

func (m *Middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := m.sessionStore.Get(r)
		if err != nil {
			response.HandleUnauthorized(w, "Invalid session")
			return
		}

		val := session.Values["user"]

		user, ok := val.(types.UserSession)

		if !ok {
			response.HandleUnauthorized(w, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user) // needs type assertion later

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
