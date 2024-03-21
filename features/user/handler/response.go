package handler

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Hp    string `json:"hp"`
	Token string `json:"token"`
}

type ProfileResponse struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Hp    string `json:"hp" form:"hp"`
}
