package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/types"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, u user.User) error
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type FileStore interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error

	Delete(ctx context.Context, objectName string) error

	// (optional) PresignedURL generates a temporary public link
	// PresignedURL(ctx context.Context, objectName string) (string, error)
}

type service struct {
	profilePicPath string
	repo           Repository
	fileStore      FileStore
}

func NewService(fs FileStore, repo Repository, profilePicPath string) *service {
	return &service{
		fileStore:      fs,
		repo:           repo,
		profilePicPath: profilePicPath,
	}
}

func (s *service) upload(ctx context.Context, f *multipart.File, h *multipart.FileHeader) (string, error) {
	var profilePicName string

	// Generate unique object name: "timestamp-filename.jpg"
	ext := filepath.Ext(h.Filename)
	profilePicName = fmt.Sprintf("%s/%d%s", s.profilePicPath, time.Now().UnixNano(), ext)

	// use context so if the user cancels the request, the upload cancels too.
	err := s.fileStore.Upload(ctx, profilePicName, *f, h.Size, h.Header.Get("Content-Type"))

	return profilePicName, err
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32) // 32 bytes provides 256 bits of entropy
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Use URLEncoding to ensure it's safe for use in URLs/headers
	return base64.URLEncoding.EncodeToString(b), nil
}

// HashToken takes a plain token string and returns the SHA-256 hash
func HashToken(token string) string {
	hash := sha256.New()

	hash.Write([]byte(token))

	return hex.EncodeToString(hash.Sum(nil))
}

// HashPassword generates a bcrypt hash of the password using a work factor
func HashPassword(password string) (string, error) {
	// GenerateFromPassword automatically adds a random "Salt"
	// Cost 10 is the default standard; 12 or 14 is better for higher security
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) // 10 = Cost
	return string(bytes), err
}

func (s *service) Register(rCtx context.Context, u *types.UserInput, parsedRefreshCookieExpiry *time.Duration, file *multipart.File, fileHeader *multipart.FileHeader) (*user.User, *string, error) {
	// check email conflict
	userExists, err := s.repo.FindByEmail(rCtx, u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// do nothing
		} else {
			return nil, nil, err
		}
	} else if userExists != nil {
		return nil, nil, errors.New("email conflict")
	}

	// upload image --> get image name (path-name/image-name.jpg)
	imagepath, err := s.upload(rCtx, file, fileHeader)

	if err != nil {
		return nil, nil, err
	}

	// refresh token
	token, err := generateRefreshToken()

	if err != nil {
		return nil, nil, err
	}
	hashedToken := HashToken(token)
	hashedPassword, err := HashPassword(u.Password)

	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	refreshTokenExpiry := now.Add(*parsedRefreshCookieExpiry)

	// insert user to db --> get user id
	err = s.repo.Create(rCtx, user.User{
		Email:              u.Email,
		Role:               "user",
		PasswordHash:       hashedPassword,
		RefreshTokenHash:   hashedToken,
		RefreshTokenExpiry: refreshTokenExpiry,
		FullName:           u.FullName,
		ProfilePicName:     imagepath,
		UpdatedAt:          now,
	})

	if err != nil {
		return nil, nil, err
	}

	// get the user row from the db
	createdUser, err := s.repo.FindByEmail(rCtx, u.Email)

	if err != nil {
		return nil, nil, err
	}

	return createdUser, &token, nil
}
