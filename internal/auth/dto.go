package auth

type RegisterRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=0"`
	FullName string `json:"full_name" validate:"required"`
}

type RegisterResponse struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	ProfilePic string `json:"profile_pic"`
}
