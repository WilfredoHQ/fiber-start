package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func accountRouter(router fiber.Router) {
	router.Get("/current", middleware.JwtAuth(), controllers.GetCurrentAccount)
	router.Post("/login", controllers.Login)
	router.Post("/recover", controllers.RecoverAccount)
	router.Post("/reset-password", controllers.ResetPassword)
}
