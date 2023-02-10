package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilfredohq/fiber-start/controllers"
	"github.com/wilfredohq/fiber-start/middleware"
)

func postRouter(router fiber.Router) {
	router.Get("", middleware.JwtAuth(), controllers.ReadPosts)
	router.Post("", middleware.JwtAuth(), controllers.CreatePost)
	router.Get("/home", middleware.JwtAuth(), controllers.ReadHomePosts)
	router.Get("/:postId", middleware.JwtAuth(), controllers.ReadPost)
	router.Delete("/:postId", middleware.JwtAuth(), controllers.DeletePost)
	router.Patch("/:postId", middleware.JwtAuth(), controllers.UpdatePost)
}
