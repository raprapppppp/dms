package routes

import (
	"dms-api/database"
	"dms-api/handlers"
	"dms-api/repository"
	"dms-api/services"

	"github.com/gofiber/fiber/v2"
)

func LoginRoute(login fiber.Router) {

	//Initializee Handler, Service, Repository
	repo := repository.LoginRepoInit(database.Database)
	service := services.LoginServicesInit(repo)
	handler := handlers.LoginHandlerInit(service)

	login.Post("/login", handler.LoginHandler)
	login.Post("/register", handler.Registerhandler)
	login.Post("/forgot-password/request", handler.ForgotPasswordRequestHandler)
	login.Post("/forgot-password/verify", handler.VerifyOTPHandler)
	//login.Post("/forgot-password/reset")
}
