package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Store struct {
	inner      sessions.Store
	cookieName string
}

func (s *Store) Get(r *http.Request) (*sessions.Session, error) {
	return s.inner.Get(r, s.cookieName)
}

func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return s.inner.Save(r, w, session)
}
