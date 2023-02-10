package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/config"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func userRouter(router fiber.Router) {
	router.Get("", middleware.JwtAuth(), controllers.GetUsers)
	if config.Config.UsersOpenRegistration {
		router.Post("", controllers.CreateUser)
	} else {
		router.Post("", middleware.JwtAuth(), controllers.CreateUser)
	}
	router.Get("/:userId", middleware.JwtAuth(), controllers.GetUser)
	router.Patch("/:userId", middleware.JwtAuth(), controllers.UpdateUser)
}
