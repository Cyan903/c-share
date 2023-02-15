package auth

import (
	"time"

	"github.com/Cyan903/c-share/pkg/config"
	"github.com/Cyan903/c-share/pkg/log"
	"github.com/golang-jwt/jwt"
)

func Create(id string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
	})

	token, err := claims.SignedString([]byte(config.Data.JWTSecret))

	if err != nil {
		log.Error.Println("Could not create JWT token -", err)
		return "", err
	}

	return token, nil
}

func VerifyToken(userToken string) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(userToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Data.JWTSecret), nil
	})

	if err != nil {
		log.Error.Println("Invalid token -", err)
		return t, err
	}

	return t, nil
}
