package user

type ProfileResponse struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsVerified bool   `json:"isVerified"`
	FullName   string `json:"fullName"`
	ProfilePic string `json:"profilePic"`
}
