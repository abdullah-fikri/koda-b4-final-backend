package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

var REFRESH_SECRET = []byte("REFRESH_SECRET_KEY")

func ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return REFRESH_SECRET, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
