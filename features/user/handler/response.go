package handler

type LoginResponse struct {
	Hp    string `json:"hp"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
