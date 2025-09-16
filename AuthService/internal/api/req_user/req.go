package requser


type RegisterRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Username string `json:"username"`
}

type RegisterResponse struct {
	Detail string `json:"detail"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Detail string `json:"detail"`
	Token string `json:"token"`
}