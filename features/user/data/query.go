package data

import (
	"21-api/features/user"
	"21-api/helper"
	"errors"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) user.Model {
	return &model{
		connection: db,
	}
}

// Register
func (m *model) Register(newData user.User) error {
	return m.connection.Create(&newData).Error
}

// Login
func (m *model) Login(email string) (user.User, error) {
	var result user.User
	if err := m.connection.Model(&User{}).Where("email = ? ", email).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

// Profil
func (m *model) Profile(id string) (user.User, error) {
	var result user.User
	err := m.connection.Where("id = ?", id).Find(&result).Error
	return result, err
}

// Berguna untuk update user
func (m *model) Update(data user.User) error {
	var selectUpdate []string
	if data.Name != "" {
		selectUpdate = append(selectUpdate, "name")
	}
	if data.Email != "" {
		selectUpdate = append(selectUpdate, "email")
	}
	if data.Password != "" {
		selectUpdate = append(selectUpdate, "password")
	}
	if data.Hp != "" {
		selectUpdate = append(selectUpdate, "hp")
	}
	if len(selectUpdate) == 0 {
		return errors.New(helper.ErrorNoRowsAffected)
	}

	if query := m.connection.Model(&data).Select(selectUpdate).Updates(&data); query.Error != nil {
		return errors.New(helper.ErrorGeneralDatabase)
	} else if query.RowsAffected == 0 {
		return errors.New(helper.ErrorNoRowsAffected)
	}
	return nil
}

// Delete User
func (m *model) Delete(id string) error {
	if query := m.connection.Where("id = ?", id).Delete(&user.User{}); query.Error != nil {
		return errors.New(helper.ErrorGeneralDatabase)
	} else if query.RowsAffected == 0 {
		return errors.New(helper.ErrorDatabaseNotFound)
	}
	return nil
}
