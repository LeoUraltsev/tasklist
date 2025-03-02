package jwt

import (
	"TaskList/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, duration time.Duration, secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["created"] = user.CreatedAt
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string, secret []byte) (jwt.MapClaims, error) {
	t, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims := t.Claims.(jwt.MapClaims)

	return claims, nil
}
