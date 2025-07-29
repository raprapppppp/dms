package handlers

import (
	"dms-api/modals"
	"dms-api/services"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type InjectLoginHandler struct {
	services services.LoginServices
}

func LoginHandlerInit(s services.LoginServices) *InjectLoginHandler {
	return &InjectLoginHandler{s}
}
//Login
func(h *InjectLoginHandler) LoginHandler(hh *fiber.Ctx) error {
	var cred modals.Login
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}
	user, mess := h.services.LoginService(cred)
	switch mess {
	case "NotFound":	
		return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Username Not Found"})
	case "NotMatch":
		return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Password Not Match"})
	case "Match":
	}
	return hh.Status(fiber.StatusOK).JSON(fiber.Map{"message" : user})
}
//Register
func(h *InjectLoginHandler) Registerhandler(hh *fiber.Ctx) error {
	var cred modals.Accounts
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	_, mess := h.services.RegisterService(cred)
	if mess == "Already Exist" {
		return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already Exist"})
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"message" : "Created"})
}
//Forgot Password
func(h *InjectLoginHandler) ForgotPasswordHandler(hh *fiber.Ctx) error {
	var email modals.Forgot
		
	if err := hh.BodyParser(&email); err != nil {
		return hh.Status(500).JSON(fiber.Map{"error": "Cannot parse JSON"}) 
	}

	otp, err := h.services.ForgotPasswordService(email)
	if err != nil {
		if errors.Is(err, services.ErrEmailNotExist){
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Email not found"})
		}
		if errors.Is(err, services.ErrOTPGenerationFailed) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate OTP. Please try again later."})
		}
	}

	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"otp": otp})
	
}