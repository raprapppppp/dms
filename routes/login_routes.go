package routes

import (
	"dms-api/database"

	"github.com/gofiber/fiber/v2"
)

func LoginRoute(login *fiber.App) {

	database.ConnectionDB()
	login.Post("/login")
}