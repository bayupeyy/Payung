package comment

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CommentController interface {
	Add() echo.HandlerFunc
	//Update() echo.HandlerFunc
	//Delete() echo.HandlerFunc
	// ShowMyTodo() echo.HandlerFunc
}

type CommentModel interface {
	InsertComment(userID string, contentBaru Comment) (Comment, error)
	UpdateComment(userID string, ID uint, data Comment) (Comment, error)
	// DeleteActivity()
	GetComment(userID string) ([]Comment, error)
}

type CommentService interface {
	AddComment(userID *jwt.Token, contentBaru Comment) (Comment, error)
	// UpdateTodo(userID *jwt.Token,
}

type Comment struct {
	Content string
}
