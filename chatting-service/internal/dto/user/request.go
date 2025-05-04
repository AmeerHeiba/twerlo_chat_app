package user

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"omitempty,min=3"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type MessageHistoryRequest struct {
	Limit  int `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int `json:"offset" validate:"omitempty,min=0"`
}
