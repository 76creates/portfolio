package server

import "github.com/gofiber/fiber/v2"

func addRootRoutes(router fiber.Router) {
	router.Static("/", "./assets")
}
