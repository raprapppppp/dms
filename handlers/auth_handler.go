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

// Login
func (h *InjectLoginHandler) LoginHandler(hh *fiber.Ctx) error {
	var cred modals.Login
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}
	user, err := h.services.LoginService(cred)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Username Not Found"})
		}
		if errors.Is(err, services.ErrNotMatch) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Password Not Match"})
		}
	}
	return hh.Status(fiber.StatusOK).JSON(fiber.Map{"message": user})
}

// Register
func (h *InjectLoginHandler) Registerhandler(hh *fiber.Ctx) error {
	var cred modals.Accounts
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	_, err := h.services.RegisterService(cred)
	if err != nil {
		if errors.Is(err, services.ErrAlreadyExist) {
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already Exist"})
		}
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Created"})
}

// Forgot Password
func (h *InjectLoginHandler) ForgotPasswordRequestHandler(hh *fiber.Ctx) error {
	var email modals.Forgot

	if err := hh.BodyParser(&email); err != nil {
		return hh.Status(500).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	mess, err := h.services.ForgotPasswordRequestService(email)
	if err != nil {
		if errors.Is(err, services.ErrEmailNotExist) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Email not found"})
		}
		if errors.Is(err, services.ErrOTPGenerationFailed) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate OTP. Please try again later."})
		}
		if errors.Is(err, services.ErrSendingOTP) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error Sending OTP."})
		}
		if errors.Is(err, services.ErrOTPRequestLimit) {
			return hh.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"message": "Please wait before requesting again."})
		}
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"Message": mess})

}
