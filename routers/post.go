package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func postRouter(router fiber.Router) {
	router.Get("", middleware.JwtAuth(), controllers.GetPosts)
	router.Post("", middleware.JwtAuth(), controllers.CreatePost)
	router.Get("/home", middleware.JwtAuth(), controllers.GetHomePosts)
	router.Get("/:postId", middleware.JwtAuth(), controllers.GetPost)
	router.Delete("/:postId", middleware.JwtAuth(), controllers.DeletePost)
	router.Patch("/:postId", middleware.JwtAuth(), controllers.UpdatePost)
}
