package middleware

import (
"net/http"
"strings"
"github.com/dgrijalva/jwt-go"
"github.com/labstack/echo/v4"
"github.com/labstack/echo/v4/middleware"
)

type JwtClaims struct {
UserId uint64 json:"user_id"
jwt.StandardClaims
}

func RequireAuth() echo.MiddlewareFunc {
return middleware.JWTWithConfig(middleware.JWTConfig{
SigningMethod: "HS256",
SigningKey: []byte("mySecretKey"),
})
}

func ExtractJwtUserId(c echo.Context) uint64 {
user := c.Get("user").(*jwt.Token)
claims := user.Claims.(*JwtClaims)
return claims.UserId
}

func TokenExtractor(header http.Header) (string, error) {
bearToken := header.Get("Authorization")
if len(bearToken) == 0 {
return "", echo.NewHTTPError(http.StatusBadRequest, "Authorization header is required")
}
strArr := strings.Split(bearToken, " ")
if len(strArr) != 2 || strings.ToLower(strArr[0]) != "bearer" {
	return "", echo.NewHTTPError(http.StatusBadRequest, "Authorization header format must be Bearer {token}")
}

return strArr[1], nil
}
