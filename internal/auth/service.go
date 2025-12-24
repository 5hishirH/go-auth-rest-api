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
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/filestore"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, u *RegisterRequest, file *multipart.File, fileHeader *multipart.FileHeader) (newUser *user.User, accessToken *string, accessTokenExipryInSec, refreshTokenExipryInSec *int, err error)
}

type service struct {
	repo      user.Repository
	fileStore filestore.FileStore
}

func (s *service) upload(ctx context.Context, f multipart.File, h *multipart.FileHeader) (string, error) {
	var profilePicName string
	defer f.Close()

	// Generate unique object name: "timestamp-filename.jpg"
	ext := filepath.Ext(h.Filename)
	profilePicName = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// use context so if the user cancels the request, the upload cancels too.
	err := s.fileStore.Upload(ctx, profilePicName, f, h.Size, h.Header.Get("Content-Type"))

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
	// 1. Create a new SHA256 hash
	hash := sha256.New()

	// 2. Write the token bytes to the hash
	hash.Write([]byte(token))

	// 3. Get the resulting bytes and encode them to a hex string
	//    The nil argument appends the result to a new slice
	return hex.EncodeToString(hash.Sum(nil))
}

// HashPassword generates a bcrypt hash of the password using a work factor
func HashPassword(password string) (string, error) {
	// GenerateFromPassword automatically adds a random "Salt"
	// Cost 10 is the default standard; 12 or 14 is better for higher security
	// but slower (which is what we want!)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) // 10 = Cost
	return string(bytes), err
}

func (s *service) Register(ctx context.Context, u *RegisterRequest, file *multipart.File, fileHeader *multipart.FileHeader) (*user.User, *string, *int, error) {
	// check email conflict
	if _, err := s.repo.FindByEmail(u.Email); err != nil {
		if err == sql.ErrNoRows {
			// the email is available.
			return nil, nil, nil, errors.New("email conflict")
		} else {
			return nil, nil, nil, err
		}
	}

	// upload image --> get image name with bucket (bucket-name/image-name.jpg)
	imagepath, err := s.upload(ctx, *file, fileHeader)

	// refresh token
	token, err := generateRefreshToken()

	if err != nil {
		return nil, nil, nil, err
	}
	hashedToken := HashToken(token)
	hashedPassword, err := HashPassword(u.Password)

	if err != nil {
		return nil, nil, nil, err
	}

	// insert user to db --> get user id
	err = s.repo.Create(ctx, user.User{
		Email:          u.Email,
		PasswordHash:   hashedPassword,
		FullName:       u.FullName,
		ProfilePicName: imagepath,
	})

	if err != nil {
		return nil, nil, nil, err
	}

	// get the user row from the db
	createdUser, err := s.repo.FindByEmail(u.Email)

	if err != nil {
		return nil, nil, nil, err
	}

	refreshTokenExipryInSec, err := s.repo.SaveRefreshToken(createdUser.Id, hashedToken)

	if err != nil {
		return nil, nil, nil, err
	}

	return createdUser, &token, &refreshTokenExipryInSec, nil
}
