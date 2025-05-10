package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3" example:"johndoe"`
	Password string `json:"password" validate:"required,min=8" example:"Password123"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@email.com"`
	Password string `json:"password" validate:"required,min=8" example:"Password123"`
}

type RefreshRequest struct {
	Token string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8" example:"Password123"`
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"NewPassword123"`
}
