package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wilfredohq/fiber-start/configs"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{AllowOrigins: strings.Join(configs.Env.BackendCorsOrigins, ",")}),
		logger.New(),
	)
}
