package routers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/constants"
	"github.com/wilfredohq/fiber-start/models"
)

func notFoundRouter(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).JSON(models.Error{Detail: constants.EndpointNotFound})
	})
}
