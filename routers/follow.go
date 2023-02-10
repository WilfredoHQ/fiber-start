package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func followRouter(router fiber.Router) {
	router.Get("", middleware.JwtAuth(), controllers.ReadFollows)
	router.Post("", middleware.JwtAuth(), controllers.CreateFollow)
	router.Get("/:followId", middleware.JwtAuth(), controllers.ReadFollow)
	router.Delete("/:followId", middleware.JwtAuth(), controllers.DeleteFollow)
}
