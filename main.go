package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/middleware"
	"github.com/wilfredohq/fiber-start/routers"
)

// @Title Start
// @Version 0.1.0
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
func main() {
	app := fiber.New()

	middleware.FiberMiddleware(app)

	routers.ApiRouter(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Listen(":" + port)
}
