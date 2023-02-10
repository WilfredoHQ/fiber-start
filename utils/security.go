package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wilfredohq/fiber-start/config"
	"golang.org/x/crypto/bcrypt"
)

func GetJwt(subject string, expirationMinutes int) (string, error) {
	expiresAt := time.Now().Add(time.Minute * time.Duration(expirationMinutes))

	claims := jwt.RegisteredClaims{Subject: subject, ExpiresAt: jwt.NewNumericDate(expiresAt)}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Config.SecretKey))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func VerifyPassword(plainPassword string, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return err
	}

	return nil
}

func GetPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
