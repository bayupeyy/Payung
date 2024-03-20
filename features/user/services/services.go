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

// Fungsi digunakan untuk memproses registrasi pengguna baru dengan menerima data pengguna baru
func (s *service) Register(newData user.User) error {
	var registerValidate user.Register
	registerValidate.Name = newData.Name
	registerValidate.Email = newData.Email
	registerValidate.Password = newData.Password
	registerValidate.Hp = newData.Hp
	err := s.v.Struct(&registerValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	newPassword, err := s.pm.HashPassword(newData.Password) //Menjadikan Password menjadi bentuk Hash
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

// Delete User
func (s *service) DeleteUser(userID string) error {
	// Lakukan logika penghapusan pengguna berdasarkan userID
	err := s.model.DeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByHP digunakan untuk mendapatkan pengguna berdasarkan nomor HP
func (s *service) GetUserByHP(hp string) (user.User, error) {
	// Memanggil model untuk mendapatkan pengguna berdasarkan nomor HP
	result, err := s.model.GetUserByHP(hp)
	if err != nil {
		return user.User{}, err
	}

	return result, nil
}

// UpdateUser digunakan untuk memperbarui pengguna berdasarkan nomor HP
func (s *service) UpdateUser(hp string, newData user.User) error {
	// Validasi data pengguna yang akan diperbarui
	err := s.v.Struct(&newData)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	// Periksa apakah pengguna yang ingin diperbarui berdasarkan nomor HP ada dalam database
	existingUser, err := s.model.GetUserByHP(hp)
	if err != nil {
		return err // Pengguna tidak ditemukan, kembalikan kesalahan
	}

	// Periksa apakah ada perubahan pada kata sandi
	if newData.Password != "" {
		// Lakukan hashing pada kata sandi baru menggunakan PasswordManager
		hashedPassword, err := s.pm.HashPassword(newData.Password)
		if err != nil {
			return err // Kembalikan kesalahan jika hashing gagal
		}
		// Gunakan kata sandi yang telah di-hash
		existingUser.Password = hashedPassword
	}

	// Perbarui data pengguna dengan informasi baru
	existingUser.Name = newData.Name
	existingUser.Email = newData.Email

	// Perbarui data pengguna di database
	err = s.model.UpdateUser(hp, existingUser)
	if err != nil {
		return err // Kembalikan kesalahan jika gagal memperbarui pengguna di database
	}

	return nil
}
