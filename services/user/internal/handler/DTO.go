package handler

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
