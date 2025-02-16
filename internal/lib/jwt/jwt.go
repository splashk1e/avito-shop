package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/splashk1e/avito-shop/internal/models"
)

func NewToken(user models.User, duration time.Duration, secret string) (string, error) {
	const op = "jwt.NewToken"
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(claims)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return tokenString, nil
}

func ParseToken(tokenString string, secret string) (username string, err error) {
	const op = "jwt.ParseToken"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", fmt.Errorf("%s: wrong jwt token", op)
		}
		return username, nil

	} else {
		return "", fmt.Errorf("%s: %w", op, err)
	}
}
