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

func (cm *model) InsertComment(PostID string, contentBaru comment.Comment) (comment.Comment, error) {
	var inputProcess = Comment{Content: contentBaru.Content, PostID: PostID}
	if err := cm.connection.Create(&inputProcess).Error; err != nil {
		return comment.Comment{}, err
	}

	return comment.Comment{Content: inputProcess.Content}, nil
}

// Fungsi untuk Edit Comment
func (cm *model) UpdateComment(postID string, commentID uint, data comment.Comment) (comment.Comment, error) {
	var qry = cm.connection.Where("id = ? AND postid = ?", commentID, postID).Updates(data)
	if err := qry.Error; err != nil {
		return comment.Comment{}, err
	}

	if qry.RowsAffected < 1 {
		return comment.Comment{}, errors.New("no data affected")
	}

	return data, nil
}

func (cm *model) GetComment(userID string) ([]comment.Comment, error) {
	var result []comment.Comment
	if err := cm.connection.Where("userid = ?", userID).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
