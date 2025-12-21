package user

type CreateRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	FullName string
}

type CreateResponse struct {
	Email          string `json:"email"`
	FullName       string `json:"full_name,omitempty"`
	ProfilePicName string `json:"profile_pic_name"`
}
