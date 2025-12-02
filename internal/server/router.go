package server

import (
	"openlist-bed/internal/server/handler"
	"openlist-bed/internal/server/middleware"

	"github.com/gofiber/fiber/v3"
)

func registerRoutes(app *fiber.App) {
	app.Get("/pic/+", handler.GetImage)

	api := app.Group("/api", middleware.Auth())

	api.Post("/upload", handler.UploadImg)
}
