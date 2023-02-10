package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/configs"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func userRouter(router fiber.Router) {
	router.Get("", middleware.JwtAuth(), controllers.ReadUsers)
	if configs.Env.UsersOpenRegistration {
		router.Post("", controllers.CreateUser)
	} else {
		router.Post("", middleware.JwtAuth(), controllers.CreateUser)
	}
	router.Get("/:userId", middleware.JwtAuth(), controllers.ReadUser)
	router.Delete("/:userId", middleware.JwtAuth(), controllers.DeleteUser)
	router.Patch("/:userId", middleware.JwtAuth(), controllers.UpdateUser)
}
