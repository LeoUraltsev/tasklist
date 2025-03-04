package jwt

import (
	"TaskList/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func NewToken(user models.User, duration time.Duration, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UID:              user.ID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration))},
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string, secret []byte) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims := t.Claims.(*CustomClaims)

	return claims, nil
}
