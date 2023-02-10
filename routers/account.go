package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func accountRouter(router fiber.Router) {
	router.Get("/current", middleware.JwtAuth(), controllers.Current)
	router.Post("/login", controllers.Login)
	router.Post("/recover-password/:email", controllers.RecoverPassword)
	router.Post("/reset-password", controllers.ResetPassword)
}
