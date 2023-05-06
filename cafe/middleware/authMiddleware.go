package middleware

import (
	"cafe/constants"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(userId int, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constants.JWTSecret))
}
