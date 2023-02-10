package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func followerRelationRouter(router fiber.Router) {
	router.Post("", middleware.JwtAuth(), controllers.CreateFollowerRelation)
	router.Get("/following/:userId", middleware.JwtAuth(), controllers.CheckFollowerRelation)
	router.Delete("/:followerRelationId", middleware.JwtAuth(), controllers.DeleteFollowerRelation)
}
