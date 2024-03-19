package services

import (
	"21-api/features/user"
	"21-api/helper"
	"21-api/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	model user.UserModel
	pm    helper.PasswordManager
	v     *validator.Validate
}

func NewService(m user.UserModel) user.UserService {
	return &service{
		model: m,
		pm:    helper.NewPasswordManager(),
		v:     validator.New(),
	}
}

// Fungsi digunakan untuk memproses registrasi pengguna baru dengan menerima data pengguna baru (nama, nomor HP, dan password) sebagai parameter input
func (s *service) Register(newData user.User) error {
	var registerValidate user.Register
	registerValidate.Nama = newData.Nama
	registerValidate.Email = newData.Email
	registerValidate.Password = newData.Password
	registerValidate.Hp = newData.Hp
	err := s.v.Struct(&registerValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	newPassword, err := s.pm.HashPassword(newData.Password)
	if err != nil {
		return errors.New(helper.ServiceGeneralError)
	}
	newData.Password = newPassword

	err = s.model.InsertUser(newData)
	if err != nil {
		return errors.New(helper.ServerGeneralError)
	}
	return nil
}

// Fungsi digunakan untuk memeriksa kredensial login pengguna yang diberikan dan mengembalikan informasi pengguna serta token akses JWT jika proses login berhasil
func (s *service) Login(loginData user.User) (user.User, string, error) {
	var loginValidate user.Login
	loginValidate.Email = loginData.Hp
	loginValidate.Password = loginData.Password
	err := s.v.Struct(&loginValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return user.User{}, "", err
	}

	dbData, err := s.model.Login(loginValidate.Email)
	if err != nil {
		return user.User{}, "", err
	}

	err = s.pm.ComparePassword(loginValidate.Password, dbData.Password)
	if err != nil {
		return user.User{}, "", errors.New(helper.UserCredentialError)
	}

	token, err := middlewares.GenerateJWT(dbData.Email)
	if err != nil {
		return user.User{}, "", errors.New(helper.ServiceGeneralError)
	}

	return dbData, token, nil
}

// Mendapatkan profil pengguna berdasarkan token akses JWT
func (s *service) Profile(token *jwt.Token) (user.User, error) {
	decodeHp := middlewares.DecodeToken(token) //Melakukan dekode token akses JWT
	result, err := s.model.GetUserByHP(decodeHp)
	if err != nil {
		return user.User{}, err
	}

	return result, nil
}
