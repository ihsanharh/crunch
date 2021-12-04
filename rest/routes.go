package controller

import (
	"github.com/XzFrosT/crunch/rest/controller/interactions"
	"github.com/XzFrosT/crunch/rest/controller/player"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/interactions", interactions.Post)
	api.Get("/player/:id", player.Index)
}
