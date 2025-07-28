package main

import (
	"dms-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	routes.LoginRoute(app)

}