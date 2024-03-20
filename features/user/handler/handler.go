package handler

import (
	"21-api/features/user"
	"21-api/helper"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service user.UserService
}

func NewUserHandler(s user.UserService) user.UserController {
	return &controller{
		service: s,
	}
}

// Register User
func (ct *controller) RegisterUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Bind data yang dikirimkan dalam request ke dalam struct User
		var input user.User
		if err := c.Bind(&input); err != nil {
			// Jika terjadi kesalahan saat binding data, tanggapi dengan kode status yang sesuai
			if strings.Contains(err.Error(), "unsupported") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirimkan tidak sesuai", nil))
		}

		// Panggil fungsi RegisterUser dari service (services.go) untuk mendaftarkan pengguna
		if err := ct.service.Register(input); err != nil {
			// Jika terjadi kesalahan saat mendaftarkan pengguna, tanggapi dengan kode status yang sesuai
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		// Tanggapi dengan status Created jika pendaftaran berhasil
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "Selamat registrasi berhasil", nil))
	}
}

func (ct *controller) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		var processData user.User
		processData.Hp = input.Hp
		processData.Password = input.Password

		result, token, err := ct.service.Login(processData)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		var responseData LoginResponse
		responseData.Hp = result.Hp
		responseData.Name = result.Name
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil login", responseData))

	}
}

func (ct *controller) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}
		result, err := ct.service.Profile(token)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

// Delete User
func (ct *controller) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Token tidak valid", nil))
		}

		// Mendapatkan klaim dari token untuk mendapatkan ID pengguna
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["id"].(string)

		// Mendapatkan ID pengguna yang akan dihapus dari parameter URL
		targetUserID := c.Param("id")

		// Memeriksa apakah pengguna mencoba menghapus akun mereka sendiri
		if userID != targetUserID {
			return c.JSON(http.StatusForbidden,
				helper.ResponseFormat(http.StatusForbidden, "Anda tidak diizinkan untuk menghapus akun pengguna lain", nil))
		}

		// Memanggil service untuk menghapus pengguna
		err := ct.service.DeleteUser(targetUserID)
		if err != nil {
			// Menangani kesalahan yang terjadi saat menghapus pengguna
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "not found") {
				code = http.StatusNotFound
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		// Mengembalikan respons JSON untuk berhasil menghapus pengguna
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil menghapus pengguna", nil))
	}
}

// Get User By HP
func (ct *controller) GetUserByHP() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan nomor HP dari parameter URL
		hp := c.Param("hp")

		// Memanggil service untuk mendapatkan pengguna berdasarkan nomor HP
		result, err := ct.service.GetUserByHP(hp)
		if err != nil {
			// Menangani kesalahan yang terjadi saat mendapatkan pengguna
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "not found") {
				code = http.StatusNotFound
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		// Mengembalikan respons JSON untuk berhasil mendapatkan pengguna berdasarkan nomor HP
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan pengguna", result))
	}
}

// Update User
func (ct *controller) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan nomor HP pengguna dari parameter URL
		hp := c.Param("hp")

		// Bind data yang dikirimkan dalam request ke dalam struct User
		var newData user.User
		if err := c.Bind(&newData); err != nil {
			// Jika terjadi kesalahan saat binding data, tanggapi dengan kode status yang sesuai
			if strings.Contains(err.Error(), "unsupported") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirimkan tidak sesuai", nil))
		}

		// Panggil fungsi UpdateUser dari service (services.go) untuk memperbarui pengguna
		if err := ct.service.UpdateUser(hp, newData); err != nil {
			// Jika terjadi kesalahan saat memperbarui pengguna, tanggapi dengan kode status yang sesuai
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "not found") {
				code = http.StatusNotFound
			} else if strings.Contains(err.Error(), "validation") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		// Tanggapi dengan status OK jika pembaruan berhasil
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil memperbarui pengguna", nil))
	}
}
