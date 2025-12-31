package auth

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/checkimage"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/response"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Register(context.Context, *RegisterRequest, *multipart.File, *multipart.FileHeader) (*user.User, error)
}

type Handler struct {
	refreshCookieName string
	refreshCookiePath string
	profileApiPrefix  string
	service           Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: *s,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// parse multipart-form

	err := r.ParseMultipartForm(10 << 20) // 10 Megabytes
	if err != nil {
		response.HandleBadRequest(w, "File too big or invalid format")
		return
	}

	// create an instance of the struct
	input := RegisterRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		FullName: r.FormValue("full_name"),
	}

	// validate
	if err := validator.New().Struct(input); err != nil {
		response.HandleValidationErrors(w, err)
		return
	}

	// get file and file header
	file, header, err := r.FormFile("profile_pic")

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

	// get created user (with refreshToken but without password & passwordHash), accessToken
	newUser, err := h.service.Register(r.Context(), &input, &file, header)
	if err != nil {
		response.HandleInternalError(w, "Error while creating user")
		return
	}

	// handle session

	// set refreshToken cookie

	// construct profile pic url
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	profilePicUrl := fmt.Sprintf("%s://%s/%s/%d", protocol, r.Host, h.profileApiPrefix, newUser.Id)

	resData := RegisterResponse{
		Id:       newUser.Id,
		Email:    newUser.Email,
		FullName: newUser.FullName,
		// Constructed: http://localhost:8080/api/v1/profile/123
		ProfilePic: profilePicUrl,
	}

	response.CreatedOne(w, "user", resData)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetVerificationEmail(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetChangePasswordEmail(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ForgetPassword(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {}
