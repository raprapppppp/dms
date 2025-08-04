package main

import (
	"dms-api/internal/database"
	"dms-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectionDB()
	app := fiber.New()

	routes.AuthRoute(app)

	app.Listen(":4000")

}
