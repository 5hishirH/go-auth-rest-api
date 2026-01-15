package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/checkimage"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/response"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/storage/session"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/types"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Register(context.Context, *types.UserInput, *time.Duration, *multipart.File, *multipart.FileHeader) (*user.User, *string, error)
	Login(context.Context, *LoginRequest, time.Duration) (*user.User, string, error)
}

type Handler struct {
	service               Service
	store                 *session.Store
	refreshCookieName     string
	refreshCookiePath     string
	refreshCookieExpiry   string
	isRefreshCookieSecure bool
	profileApiPrefix      string
}

func NewHandler(s Service, store *session.Store, refreshCookieName, refreshCookiePath, refreshCookieExpiry string, isRefreshCookieSecure bool, profileApiPrefix string) *Handler {
	return &Handler{
		service:               s,
		store:                 store,
		refreshCookieName:     refreshCookieName,
		refreshCookiePath:     refreshCookiePath,
		refreshCookieExpiry:   refreshCookieExpiry,
		isRefreshCookieSecure: isRefreshCookieSecure,
		profileApiPrefix:      profileApiPrefix,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r)

	if err != nil {
		response.HandleInternalError(w, "Error while initiating session")
		return
	}

	// parse multipart-form
	err = r.ParseMultipartForm(10 << 20) // 10 Megabytes
	if err != nil {
		response.HandleBadRequest(w, "File too big or invalid format")
		return
	}

	// create an instance of the struct
	req := RegisterRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		FullName: r.FormValue("fullName"),
	}

	// validate
	if err := validator.New().Struct(req); err != nil {
		response.HandleValidationErrors(w, err)
		return
	}

	// get file and file header
	file, header, err := r.FormFile("profilePic")

	// what type of errors does r.FormFile through?
	if err != nil {
		response.HandleBadRequest(w, "Profile picture is required")
		return
	}

	defer file.Close()

	// check file if image
	if err = checkimage.CheckImage(file); err != nil {
		response.HandleBadRequest(w, "Profile picture file should be of type jpg or png")
		return
	}

	// no benefit in separte function as error to be handled in each usage separately
	parsedRefreshCookieExpiry, err := time.ParseDuration(h.refreshCookieExpiry)

	if err != nil {
		fmt.Printf("refresh token expiry: %s", h.refreshCookieExpiry)
		fmt.Printf("%s", err.Error())
		response.HandleInternalError(w, "Error in parsing refresh cookie duration")
		return
	}

	userInput := types.UserInput{
		Email:      req.Email,
		Password:   req.Password,
		Role:       "user",
		IsVerified: false,
		FullName:   req.FullName,
	}

	// get created user (with refreshToken but without password & passwordHash), accessToken
	newUser, refreshToken, err := h.service.Register(r.Context(), &userInput, &parsedRefreshCookieExpiry, &file, header)
	if err != nil {
		fmt.Printf("from service: %s\n", err.Error())
		if err.Error() == "email conflict" {
			response.HandleConflict(w, "Email already exists")
			return
		} else {

			response.HandleInternalError(w, "Error while creating user")
			return
		}
	}

	// handle session
	userSession := types.UserSession{
		UserID: newUser.Id,
		Role:   newUser.Role,
	}

	session.Values["user"] = userSession
	session.Save(r, w)

	refreshCookieExpiryInSec := int(parsedRefreshCookieExpiry.Seconds())
	// set refreshToken cookie
	GenerateCookieResponse(w, h.refreshCookieName, h.refreshCookiePath, *refreshToken, refreshCookieExpiryInSec, h.isRefreshCookieSecure)

	// construct profile pic url
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	profilePicUrl := fmt.Sprintf("%s://%s/%s/%s", protocol, r.Host, h.profileApiPrefix, "profile-pic")

	resData := RegisterResponse{
		Id:         newUser.Id,
		Email:      newUser.Email,
		Role:       newUser.Role,
		IsVerified: newUser.IsVerified,
		FullName:   newUser.FullName,
		// Constructed: http://localhost:8080/api/v1/profile/123
		ProfilePic: profilePicUrl,
	}

	response.CreatedOne(w, "user", resData)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r)

	if err != nil {
		response.HandleInternalError(w, "Error while initiating session")
		return
	}

	var loginCredentials LoginRequest

	err = json.NewDecoder(r.Body).Decode(&loginCredentials)

	if errors.Is(err, io.EOF) {
		response.HandleBadRequest(w, "empty body")
		return
	}

	if err != nil {
		response.HandleBadRequest(w, err.Error())
		return
	}

	if err := validator.New().Struct(loginCredentials); err != nil {
		response.HandleValidationErrors(w, err)
		return
	}

	parsedRefreshCookieExpiry, err := time.ParseDuration(h.refreshCookieExpiry)

	if err != nil {
		fmt.Printf("refresh token expiry: %s", h.refreshCookieExpiry)
		fmt.Printf("%s", err.Error())
		response.HandleInternalError(w, "Error in parsing refresh cookie duration")
		return
	}

	user, refreshToken, err := h.service.Login(r.Context(), &loginCredentials, parsedRefreshCookieExpiry)

	if err != nil {
		if err.Error() == "invalid credentials" {
			response.HandleUnauthorized(w, err.Error())
			return
		}

		response.HandleInternalError(w, err.Error())
		return
	}

	// handle session
	userSession := types.UserSession{
		UserID: user.Id,
		Role:   user.Role,
	}

	session.Values["user"] = userSession
	session.Save(r, w)

	refreshCookieExpiryInSec := int(parsedRefreshCookieExpiry.Seconds())

	GenerateCookieResponse(w, h.refreshCookieName, h.refreshCookiePath, refreshToken, refreshCookieExpiryInSec, h.isRefreshCookieSecure)

	// construct profile pic url
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	profilePicUrl := fmt.Sprintf("%s://%s/%s/%s", protocol, r.Host, h.profileApiPrefix, "profile-pic")

	resData := RegisterResponse{
		Id:         user.Id,
		Email:      user.Email,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		FullName:   user.FullName,
		// Constructed: http://localhost:8080/api/v1/profile/123
		ProfilePic: profilePicUrl,
	}

	response.CreatedOne(w, "user", resData)
}

func (h *Handler) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetVerificationEmail(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetChangePasswordEmail(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ForgetPassword(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// revoke session
	session, err := h.store.Get(r)

	if err != nil {
		response.HandleInternalError(w, "Error while initiating session")
		return
	}

	session.Options.MaxAge = -1
	session.Save(r, w)

	// revoke refresh token
	GenerateClearCookieResponse(w, h.refreshCookieName, h.refreshCookiePath)
	response.NoContent(w)
}
