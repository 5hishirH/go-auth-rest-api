package auth

import (
	"context"
	"mime/multipart"

	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
)

type Service interface {
	Register(ctx context.Context, u *RegisterRequest, file *multipart.File, fileHeader *multipart.FileHeader) (newUser *user.User, accessToken *string, accessTokenExipryInSec, refreshTokenExipryInSec *int, err error)
}
