package router

import (
	"github.com/dixiedream/authenticator/handler"
	"github.com/dixiedream/authenticator/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", logger.New(), middleware.ValidateIP)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/", handler.Login)
	auth.Delete("/", handler.Logout)
	auth.Post("/refresh", handler.Refresh)

    // Token
    jwt := api.Group("/jwt", middleware.Auth())
    jwt.Post("/", handler.CheckToken)

	// Server
	server := api.Group("/server", middleware.Auth())
	server.Post("/", middleware.IsAdmin, handler.CreateServer)
}
