package auth

type AuthResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"` // seconds
	TokenType    string `json:"token_type" example:"Bearer"`
	UserID       uint   `json:"user_id" example:"1"`
	Username     string `json:"username" example:"johndoe"`
	Email        string `json:"email" example:"john@example.com"`
}
