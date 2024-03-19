package data

import (
	"21-api/features/comment"
	"errors"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) comment.CommentModel {
	return &model{
		connection: db,
	}
}

func (cm *model) InsertComment(PostID uint, commentBaru comment.CommentModel) (comment.CommentModel, error) {
	var inputProcess = Comment{Content: commentBaru.Comment, PostID: PostID}
	if err := cm.connection.Create(&inputProcess).Error; err != nil {
		return comment.CommentModel{}, err
	}

	return comment.CommentModel{Content: inputProcess.Content}, nil
}

// Fungsi untuk Edit Comment
func (cm *model) UpdateComment(PostID uint, data comment.CommentModel) (comment.CommentModel, error) {
	var qry = cm.connection.Where("postid = ?", PostID).Updates(data)
	if err := qry.Error; err != nil {
		return comment.CommentModel{}, err
	}

	if qry.RowsAffected < 1 {
		return comment.CommentModel{}, errors.New("no data affected")
	}

	return data, nil
}

func (cm *model) GetCommentByOwner(PostID uint) ([]comment.CommentModel, error) {
	var result []comment.CommentModel
	if err := cm.connection.Where("postid = ?", PostID).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
