package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/wilfredohq/fiber-start/docs"
)

func swaggerRouter(router fiber.Router) {
	router.Get("/*", swagger.HandlerDefault)
}
