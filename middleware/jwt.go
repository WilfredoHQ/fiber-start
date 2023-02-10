package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wilfredohq/fiber-start/config"
	"github.com/wilfredohq/fiber-start/constants"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JwtAuth() func(*fiber.Ctx) error {
	config := jwtware.Config{
		SigningKey:     []byte(config.Config.SecretKey),
		ContextKey:     "jwt",
		Claims:         &jwt.RegisteredClaims{},
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
	}

	return jwtware.New(config)
}

func jwtSuccess(c *fiber.Ctx) error {
	claims := utils.GetLocalJwtClaims(c)

	userID, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: constants.InvalidJwt})
	}

	c.Locals("userId", userID)

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: constants.InvalidJwt})
}
