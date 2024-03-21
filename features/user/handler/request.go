package handler

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	Nama     string
	Email    string
	Password string
	Hp       string
}
