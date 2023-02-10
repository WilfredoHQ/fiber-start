package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wilfredohq/fiber-start/config"
)

func GetLocalJwtClaims(c *fiber.Ctx) *jwt.RegisteredClaims {
	token := c.Locals("jwt").(*jwt.Token)
	return token.Claims.(*jwt.RegisteredClaims)
}

func GetJwtClaims(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.SecretKey), nil
	})
	if err != nil {
		return &jwt.RegisteredClaims{}, err
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}
