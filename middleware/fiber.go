package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wilfredohq/fiber-start/config"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{AllowOrigins: config.Config.BackendCorsOrigins}),
		logger.New(),
	)
}
