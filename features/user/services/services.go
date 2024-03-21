package services

import (
	"21-api/features/user"
	"21-api/helper"
	"21-api/middlewares"
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	model user.Model
	pm    helper.PasswordManager
	v     *validator.Validate
}

func NewService(m user.Model) user.Service {
	return &service{
		model: m,
		pm:    helper.NewPasswordManager(),
		v:     validator.New(),
	}
}

// Fungsi digunakan untuk memproses registrasi pengguna baru dengan menerima data pengguna baru
func (s *service) Register(newData user.Register) error {
	// Check Validate
	var validate user.Register
	validate.Name = newData.Name
	validate.Email = newData.Email
	validate.Password = newData.Password
	validate.Hp = newData.Hp

	err := s.v.Struct(&validate)
	if err != nil {
		return errors.New(helper.ErrorInvalidValidate)
	}

	// Hashing Password
	newPassword, err := s.pm.HashPassword(newData.Password)
	if err != nil {
		return errors.New(helper.ErrorGeneralServer)
	}

	user_data := user.User{
		Name:  newData.Name,
		Email: newData.Email,

		Password: newPassword,
		Hp:       newData.Hp,
	}

	// Do Register
	err = s.model.Register(user_data)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			return errors.New(mysqlErr.Message)
		}
		return errors.New(helper.ErrorGeneralServer)
	}

	return nil
}

// Fungsi digunakan untuk memeriksa kredensial login pengguna yang diberikan dan mengembalikan informasi pengguna serta token akses JWT jika proses login berhasil
func (s *service) Login(loginData user.User) (user.LoginResponse, error) {
	// Do Login & Get Password
	dbData, err := s.model.Login(loginData.Email)
	if err != nil {
		return user.LoginResponse{}, errors.New(helper.ErrorDatabaseNotFound)
	}

	// Compare Password
	if err := s.pm.ComparePassword(loginData.Password, dbData.Password); err != nil {
		return user.LoginResponse{}, errors.New(helper.ErrorUserCredential)
	}

	// Convert dbData.ID to int
	id, err := strconv.Atoi(dbData.ID)
	if err != nil {
		return user.LoginResponse{}, err // Handle error if conversion fails
	}

	// Create Token
	token, err := middlewares.GenerateJWT(strconv.Itoa(id))
	if err != nil {
		return user.LoginResponse{}, errors.New(helper.ErrorGeneralServer)
	}

	// Finished
	var result user.LoginResponse
	result.Name = dbData.Name
	result.Email = dbData.Email
	result.Hp = dbData.Hp
	result.Token = token
	return result, nil
}

// Mendapatkan profil pengguna berdasarkan token akses JWT
func (s *service) Profile(token *jwt.Token) (user.User, error) {
	// Get ID From Token
	decodeID := middlewares.DecodeToken(token)

	// Get Profile
	result, err := s.model.Profile(decodeID)
	if err != nil {
		return user.User{}, err
	}

	// Finished
	return result, nil
}

// UpdateUser digunakan untuk memperbarui pengguna berdasarkan nomor HP
func (s *service) Update(token *jwt.Token, updateData user.User) error {
	// Get ID From Token
	decodeID := middlewares.DecodeToken(token)

	// Check Validate Password & Others
	var validate user.Update
	validate.Name = updateData.Name
	validate.Email = updateData.Email

	validate.Password = updateData.Password
	validate.Hp = updateData.Hp
	err := s.v.Struct(&validate)
	if err != nil {
		if strings.Contains(err.Error(), "Name") {
			updateData.Name = ""
		}

		if strings.Contains(err.Error(), "Email") {
			updateData.Email = ""
		}
		if strings.Contains(err.Error(), "Password") {
			updateData.Password = ""
		}
		if strings.Contains(err.Error(), "Hp") {
			updateData.Password = ""
		}
	}

	// Convert id to string
	updateData.ID = decodeID

	// Hashing Password
	if updateData.Password != "" {
		newPassword, err := s.pm.HashPassword(updateData.Password)
		if err != nil {
			return errors.New(helper.ErrorGeneralServer)
		}
		updateData.Password = newPassword
	}

	// Update Data
	if err := s.model.Update(updateData); err != nil {
		return err
	}

	// Finished
	return nil
}

func (s *service) Delete(token *jwt.Token) error {
	// Get ID From Token
	decodeID := middlewares.DecodeToken(token)

	// Delete Date
	if err := s.model.Delete(decodeID); err != nil {
		return err
	}

	// Finished
	return nil
}
