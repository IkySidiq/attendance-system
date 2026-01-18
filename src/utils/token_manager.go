package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var AccessTokenKey = []byte(os.Getenv("ACCESS_TOKEN_KEY"))
var RefreshTokenKey = []byte(os.Getenv("REFRESH_TOKEN_KEY"))

func GenerateAccessToken(payload map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": payload["user_id"],
		"exp":     time.Now().Add(time.Second * 1800).Unix(), // ACCESS_TOKEN_AGE
	})
	return token.SignedString(AccessTokenKey)
}

func GenerateRefreshToken(payload map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": payload["user_id"],
	})
	return token.SignedString(RefreshTokenKey)
}

func VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return AccessTokenKey, nil
	})
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return RefreshTokenKey, nil
	})
}