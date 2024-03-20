package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	RegisterUser() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	GetUserByHP() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
}

type UserService interface {
	Register(newData User) error
	Login(loginData User) (User, string, error)
	Profile(token *jwt.Token) (User, error)
	DeleteUser(userID string) error
	GetUserByHP(hp string) (User, error)
	UpdateUser(hp string, newData User) error
}

type UserModel interface {
	InsertUser(newData User) error
	UpdateUser(hp string, data User) error
	Login(hp string) (User, error)
	GetUserByHP(hp string) (User, error)
	DeleteUser(userID string) error
}

type User struct {
	Name     string
	Email    string
	Password string
	Hp       string
}

type Login struct {
	Email    string
	Password string `validate:"required,alphanum,min=8"`
}

type Register struct {
	Name     string `validate:"required,alpha"`
	Email    string
	Password string `validate:"required,alphanum,min=8"`
	Hp       string `validate:"required,min=10,max=13,numeric"`
}
