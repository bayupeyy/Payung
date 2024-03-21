package helper

import "net/http"

const (
	// User Input
	ErrorUserInput       = "data yang dikirim tidak sesuai"
	ErrorUserInputFormat = "format data tidak didukung"
	ErrorUserCredential  = "email atau kata sandi salah"
	ErrorInvalidValidate = "validasi tidak valid"

	// Server
	ErrorGeneralServer = "terjadi kesalahan dalam proses server"

	// Database
	ErrorGeneralDatabase  = "terdapat masalah dengan database"
	ErrorNoRowsAffected   = "tidak ada perubahan pada database"
	ErrorDatabaseNotFound = "data tidak ditemukan pada database"
)

// ErrorCode maps an error to an HTTP status code.
func ErrorCode(e error) int {
	if e == nil {
		return http.StatusOK // 200
	}
	switch e.Error() {
	// User Input
	case ErrorUserInput:
		return http.StatusBadRequest // 400
	case ErrorUserCredential:
		return http.StatusUnauthorized // 401
	case ErrorInvalidValidate:
		return http.StatusBadRequest // 400

	// Server
	case ErrorGeneralServer:
		return http.StatusInternalServerError // 500

	// Database
	case ErrorGeneralDatabase:
		return http.StatusInternalServerError // 500
	case ErrorDatabaseNotFound:
		return http.StatusNotFound // 404

	// Default
	default:
		return http.StatusBadRequest // 400
	}
}
