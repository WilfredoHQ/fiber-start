package routers

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRouter(app *fiber.App) {
	swaggerRouter(app.Group("/swagger"))

	prefix := "/api/v1"
	accountRouter(app.Group(prefix + "/account"))
	followRouter(app.Group(prefix + "/follows"))
	postRouter(app.Group(prefix + "/posts"))
	userRouter(app.Group(prefix + "/users"))

	notFoundRouter(app)
}
