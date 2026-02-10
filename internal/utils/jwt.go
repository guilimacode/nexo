package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func getSecretKey() []byte {
	if len(secretKey) == 0 {
		key := os.Getenv("SECRET_KEY")
		secretKey = []byte(key)
	}
	return secretKey
}

func GenerateToken(email, role string, userID, OrganizationID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"org":   OrganizationID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 8).Unix(),
	})

	return token.SignedString(getSecretKey())
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
