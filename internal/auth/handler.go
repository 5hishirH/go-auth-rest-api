package auth

import (
	"database/sql"
	"fmt"
	"net/http"

	checkimage "github.com/5hishirH/go-auth-rest-api.git/internal/shared/utils/check-image"
	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/utils/response"
	"github.com/5hishirH/go-auth-rest-api.git/internal/user"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	accessCookieName  string
	accessCookiePath  string
	refreshCookieName string
	refreshCookiePath string
	profileApiPrefix  string
	service           Service
	userRepo          user.Repository
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
		response.WriteJSON(w, http.StatusBadRequest, *response.GeneralError("File too big or invalid format"))
		return
	}

	// create an instance of the struct
	req := RegisterRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		FullName: r.FormValue("full_name"),
	}

	// validate
	if err := validator.New().Struct(req); err != nil {
		response.HandleValidationErrors(w, err)
		return
	}

	// check duplicate email
	if _, err := h.userRepo.FindByEmail(req.Email); err != nil {
		if err == sql.ErrNoRows {
			// the email is available.
			response.HandleConflict(w, "An account with this email already exists")
			return
		} else {
			//  Any other error is a real DB issue (connection died, syntax error, etc)
			response.HandleInternalError(w, "Unknown error while checking email")
			return
		}
	}

	// get file and file header
	file, header, err := r.FormFile("profile_pic")

	if err != nil {
		response.HandleBadRequest(w, "Profile picture is required")
		return
	}

	defer file.Close()

	// check file if image
	if err = checkimage.CheckImage(file); err != nil {
		response.HandleBadRequest(w, "profile picture file should be of type jpg or png")
		return
	}

	// get created user (with refreshToken but without password & passwordHash), accessToken
	newUser, accessToken, accessCookieExpiry, refreshTokenExipry, err := h.service.Register(r.Context(), &req, &file, header)
	if err != nil {
		response.HandleInternalError(w, "Error while creating user")
		return
	}

	// set accessToken & refreshToken cookie
	h.SetAuthCookies(w, *accessToken, *accessCookieExpiry, newUser.RefreshToken, *refreshTokenExipry)

	// construct profile pic url
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	profilePicUrl := fmt.Sprintf("%s://%s/%s/%s", protocol, r.Host, h.profileApiPrefix, newUser.Id)

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
