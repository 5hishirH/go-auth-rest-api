package auth

type RegisterRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=0"`
	FullName string `json:"fullName" validate:"required"`
}

type RegisterResponse struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsVerified bool   `json:"isVerified"`
	FullName   string `json:"fullName"`
	ProfilePic string `json:"profilePic"`
}
