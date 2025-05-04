package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RefreshRequest struct {
	Token string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
