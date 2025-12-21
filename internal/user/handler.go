package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/5hishirH/go-auth-rest-api.git/internal/shared/utils/response"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("creating a user")

	// parse multipart-form
	// create an instance of the struct
	// validate the struct

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too big or invalid format", http.StatusBadRequest)
		return
	}

	req := CreateRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		FullName: r.FormValue("full_name"),
	}

	if err := validator.New().Struct(req); err != nil {
		response.HandleValidationErrors(w, err)
		return
	}

	file, header, err := r.FormFile("profile_pic")
	var profilePicName string

	if err == nil {
		defer file.Close()

		// Generate unique object name: "timestamp-filename.jpg"
		ext := filepath.Ext(header.Filename)
		profilePicName = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		// 4. Upload to MinIO
		// We use r.Context() so if the user cancels the request, the upload cancels too.
		err = h.fileStore.Upload(r.Context(), profilePicName, file, header.Size, header.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, "Failed to upload image", http.StatusInternalServerError)
			return
		}
	}

	response.CreatedOne(w, "user", CreateResponse{
		Email:          req.Email,
		FullName:       req.FullName,
		ProfilePicName: profilePicName,
	})
}
