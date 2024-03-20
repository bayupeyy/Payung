package data

import (
	"21-api/features/user"
	"errors"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) user.UserModel {
	return &model{
		connection: db,
	}
}

func (m *model) InsertUser(newData user.User) error {
	err := m.connection.Create(&newData).Error

	if err != nil {
		return errors.New("terjadi masalah pada database")
	}

	return nil
}

func (m *model) cekUser(hp string) bool {
	var data User
	if err := m.connection.Where("hp = ?", hp).First(&data).Error; err != nil {
		return false
	}
	return true
}

// Berguna untuk update user
func (m *model) UpdateUser(hp string, newData user.User) error {
	// Mencari pengguna yang ingin diperbarui berdasarkan nomor HP
	var existingUser user.User
	if err := m.connection.Where("hp = ?", hp).First(&existingUser).Error; err != nil {
		return err
	}

	// Memperbarui informasi pengguna yang ditemukan dengan data baru
	existingUser.Name = newData.Name
	existingUser.Email = newData.Email
	existingUser.Password = newData.Password

	// Menyimpan perubahan ke dalam database dengan menambahkan kondisi WHERE
	if err := m.connection.Where("hp = ?", hp).Save(&existingUser).Error; err != nil {
		return err
	}

	return nil
}

// Get User By Hp
func (m *model) GetUserByHP(hp string) (user.User, error) {
	var result user.User
	if err := m.connection.Model(&user.User{}).Where("hp = ?", hp).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

// Login
func (m *model) Login(email string) (user.User, error) {
	var result user.User
	if err := m.connection.Model(&User{}).Where("email = ? ", email).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

// Delete User
func (m *model) DeleteUser(userID string) error {
	// Membuat objek user.User dengan ID yang akan dihapus
	userToDelete := user.User{Hp: userID}

	// Menghapus pengguna dari database berdasarkan ID
	if err := m.connection.Delete(&userToDelete).Error; err != nil {
		return err
	}

	return nil
}
