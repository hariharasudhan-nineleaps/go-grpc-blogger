package utils

import "github.com/golang-jwt/jwt/v5"

func GenerateToken(claims *jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}
