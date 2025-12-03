package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

type UserPayload struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// generate token
func GenerateAccessToken(id int, role string) string {
	claims := UserPayload{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(accessSecret)
	return signed
}

// refresh token
func GenerateRefreshToken(id int, role string) string {
	claims := UserPayload{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(refreshSecret)
	return signed
}


func GenerateTokens(id int, role string) map[string]string {
	return map[string]string{
		"access_token":  GenerateAccessToken(id, role),
		"refresh_token": GenerateRefreshToken(id, role),
	}
}


func VerifyAccessToken(token string) (UserPayload, error) {
	parsed, err := jwt.ParseWithClaims(token, &UserPayload{}, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if err != nil {
		return UserPayload{}, err
	}
	return *(parsed.Claims.(*UserPayload)), nil
}

func VerifyRefreshToken(token string) (UserPayload, error) {
	parsed, err := jwt.ParseWithClaims(token, &UserPayload{}, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return UserPayload{}, err
	}
	return *(parsed.Claims.(*UserPayload)), nil
}
