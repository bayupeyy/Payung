package activity

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CommentController interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	// ShowMyTodo() echo.HandlerFunc
}

type CommentModel interface {
	InsertComment(pemilik string, kegiatanBaru Comment) (Comment, error)
	UpdateComment(pemilik string, activityID uint, data Comment) (Comment, error)
	// DeleteActivity()
	GetCommentByOwner(pemilik string) ([]Comment, error)
}

type ActivityService interface {
	AddComment(pemilik *jwt.Token, kegiatanBaru Comment) (Comment, error)
	// UpdateTodo(pemilik *jwt.Token, todoID string, data Todo) (Todo, error)
}

type Comment struct {
	Kegiatan string
}
