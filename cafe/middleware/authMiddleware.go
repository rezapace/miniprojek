package middleware

// Import library yang dibutuhkan
import (
	"cafe/constants"
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO:  Fungsi CreateToken digunakan untuk membuat token baru
func CreateToken(userId int, name string, userRole string) (string, error) {
	// Deklarasi variabel claims sebagai jwt.MapClaims
	claims := jwt.MapClaims{}
	// Menambahkan user_id dan name ke dalam claims
	claims["user_id"] = userId
	claims["name"] = name
	claims["user_role"] = userRole
	// Menambahkan expiry date ke dalam claims
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Membuat token baru dengan metode SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Mengembalikan token yang telah ditandatangani dengan secret key yang ada di constants.JWTSecret
	return token.SignedString([]byte(constants.JWTSecret))
}
