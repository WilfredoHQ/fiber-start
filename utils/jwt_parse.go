package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wilfredohq/fiber-start/configs"
)

func GetLocalJwtClaims(c *fiber.Ctx) *jwt.RegisteredClaims {
	token := c.Locals("jwt").(*jwt.Token)
	return token.Claims.(*jwt.RegisteredClaims)
}

func GetJwtClaims(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(tokenString *jwt.Token) (interface{}, error) {
		return []byte(configs.Env.SecretKey), nil
	})
	if err != nil {
		return token.Claims.(*jwt.RegisteredClaims), err
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}
