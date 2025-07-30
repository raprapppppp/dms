package routes

import (
	"dms-api/internal/database"
	"dms-api/internal/handlers"
	"dms-api/internal/repository"
	"dms-api/internal/services"

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
