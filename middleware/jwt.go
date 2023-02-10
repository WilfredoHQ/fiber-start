package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wilfredohq/fiber-start/configs"
	"github.com/wilfredohq/fiber-start/models"
	"github.com/wilfredohq/fiber-start/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JwtAuth() func(*fiber.Ctx) error {
	config := jwtware.Config{
		SigningKey:     []byte(configs.Env.SecretKey),
		ContextKey:     "jwt",
		Claims:         &jwt.RegisteredClaims{},
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
	}

	return jwtware.New(config)
}

func jwtSuccess(c *fiber.Ctx) error {
	claims := utils.GetLocalJwtClaims(c)

	if !primitive.IsValidObjectID(claims.Subject) {
		return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: "Failed to validate credentials"})
	}

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(http.StatusBadRequest).JSON(models.Error{Detail: err.Error()})
	}

	return c.Status(http.StatusUnauthorized).JSON(models.Error{Detail: "Failed to validate credentials"})
}
