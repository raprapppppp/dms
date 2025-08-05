package routes

import (
	"dms-api/internal/database"
	"dms-api/internal/handlers"
	"dms-api/internal/middleware"
	"dms-api/internal/repository"
	"dms-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(auth fiber.Router) {

	//Initializee Handler, Service, Repository
	repo := repository.AuthRepoInit(database.Database)
	service := services.AuthServicesInit(repo)
	handler := handlers.AuthHandlerInit(service)

	auth.Post("/login", handler.LoginHandler)
	auth.Post("/register", handler.Registerhandler)
	auth.Post("/forgot-password/request", handler.ForgotPasswordRequestHandler)
	auth.Post("/forgot-password/verify", handler.VerifyOTPHandler)
	auth.Post("/forgot-password/reset", middleware.PasswordResetAuthMiddleware, handler.PasswordResetHandler)

	//auth.Post("/encryption", handlers.EncryptionHandler)
	//auth.Post("/decryption", handlers.DecryptionHandler)
}
