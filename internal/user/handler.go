package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/response"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/types"
)

type Repository interface {
	FindById(ctx context.Context, id int64) (*User, error)
}

type Handler struct {
	requireAuth      func(http.HandlerFunc) http.HandlerFunc
	repo             Repository
	profileApiPrefix string
}

func NewHandler(requireAuth func(http.HandlerFunc) http.HandlerFunc, repo Repository, profileApiPrefix string) *Handler {
	return &Handler{
		requireAuth:      requireAuth,
		repo:             repo,
		profileApiPrefix: profileApiPrefix,
	}
}

var userKey string = "user"

func GetUserFromContext(ctx context.Context) (types.UserSession, bool) {
	// Retrieve the value
	val := ctx.Value(userKey)

	// Type assertion: Check if it's actually a User struct
	user, ok := val.(types.UserSession)
	return user, ok
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	u, ok := GetUserFromContext(r.Context())

	if !ok {
		response.HandleInternalError(w, "Error retriving user from session")
		return
	}

	user, err := h.repo.FindById(r.Context(), u.UserID)

	if err != nil {
		response.HandleInternalError(w, "Error retriving user from db")
		return
	}

	// construct profile pic url
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	profilePicUrl := fmt.Sprintf("%s://%s/%s/%s", protocol, r.Host, h.profileApiPrefix, "profile-pic")

	resData := ProfileResponse{
		Id:         user.Id,
		Email:      user.Email,
		FullName:   user.FullName,
		ProfilePic: profilePicUrl,
		Role:       user.Role,
		IsVerified: user.IsVerified,
	}

	response.Retrived(w, "user", resData)
}
