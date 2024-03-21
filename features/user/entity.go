package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Controller interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	Register(newData Register) error
	Login(loginData User) (LoginResponse, error)
	Profile(token *jwt.Token) (User, error)
	Update(token *jwt.Token, updateData User) error
	Delete(token *jwt.Token) error
}

type Model interface {
	Register(newData User) error
	Login(email string) (User, error)
	Profile(id string) (User, error)
	Update(data User) error
	Delete(id string) error
}

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string `json:"password" form:"password" validate:"required"`
	Hp       string `gorm:"type:varchar(13);uniqueIndex;not null" json:"hp" form:"hp" validate:"required,max=13,min=10"`
}
type Login struct {
	Email    string `validate:"required"`
	Password string `validate:"required,alphanum,min=8"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Hp    string `json:"hp"`
	Token string `json:"token"`
}
type Register struct {
	Name     string `validate:"required,alpha"`
	Email    string
	Password string `validate:"required,alphanum,min=8"`
	Hp       string `validate:"required,min=10,max=13,numeric"`
}

type Update struct {
	Name     string `validate:"required,min=5"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Hp       string `validate:"required,number,min=11,max=14"`
}
