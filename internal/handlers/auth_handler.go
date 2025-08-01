package handlers

import (
	"dms-api/internal/modals"
	"dms-api/internal/services"

	//"dms-api/utils/cryptography/decrypt"
	//"dms-api/utils/cryptography/encrypt"
	"dms-api/utils/customerror"
	"dms-api/utils/jwtgenerator"
	"errors"
	"time"

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
		return hh.Status(500).JSON(fiber.Map{"message": customerror.ParseError})
	}
	user, err := h.services.LoginService(cred)
	if err != nil {
		if errors.Is(err, customerror.ErrNotFound) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		if errors.Is(err, customerror.ErrNotMatch) {
			return hh.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Password Incorrect"})
		}
	}
	return hh.Status(fiber.StatusOK).JSON(fiber.Map{"message": user})
}

// Register
func (h *InjectLoginHandler) Registerhandler(hh *fiber.Ctx) error {
	var cred modals.Accounts
	if err := hh.BodyParser(&cred); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": customerror.ParseError})
	}
	_, err := h.services.RegisterService(cred)
	if err != nil {
		if errors.Is(err, customerror.ErrAlreadyExist) {
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": customerror.ErrAlreadyExist})
		}
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Created"})
}

// Forgot Password
func (h *InjectLoginHandler) ForgotPasswordRequestHandler(hh *fiber.Ctx) error {
	var email modals.Forgot

	if err := hh.BodyParser(&email); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": customerror.ParseError})
	}

	mess, err := h.services.ForgotPasswordRequestService(email)
	if err != nil {
		if errors.Is(err, customerror.ErrEmailNotExist) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Email does not exist"})
		}
		if errors.Is(err, customerror.ErrOTPGenerationFailed) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "OTP generate failed"})
		}
		if errors.Is(err, customerror.ErrSendingOTP) {
			return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed sending OTP"})
		}
		if errors.Is(err, customerror.ErrOTPRequestLimit) {
			return hh.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"message": "To many attempts"})
		}
	}
	return hh.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": mess})
}

// Verify OTP
func (h *InjectLoginHandler) VerifyOTPHandler(hh *fiber.Ctx) error {
	var otp modals.VerifyOTP

	if err := hh.BodyParser(&otp); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": customerror.ParseError})
	}
	id, err := h.services.VerifyOTPService(otp)
	if err != nil {
		if errors.Is(err, customerror.ErrNotFound) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Email not found"})
		}
		if errors.Is(err, customerror.ErrNoLatestOtp) {
			return hh.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No Password reset request"})
		}
		if errors.Is(err, customerror.ErrOTPAlreadyUsed) {
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "OTP already used"})
		}
		if errors.Is(err, customerror.ErrOTPExpired) {
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "OTP Expired"})
		}
		if errors.Is(err, customerror.ErrNotMatch) {
			return hh.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "OTP not match"})
		}
		if errors.Is(err, customerror.ErrOTPUpdateFailed) {
			return hh.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to updated OTP status"})
		}
	}
	//Generate jwt token
	token, err := jwtgenerator.GenerateOTPToken(id)
	if err != nil {
		return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sign token",
		})
	}
	//Creating Cookies struct may other way setcookie
	hh.Cookie(&fiber.Cookie{
		Name:     "otpToken",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	return hh.Status(fiber.StatusOK).JSON(fiber.Map{"message": token})
}

//Password Reset
func(h *InjectLoginHandler) PasswordResetHandler(hh *fiber.Ctx) error {
	var newPassword modals.NewPassword

	if err := hh.BodyParser(&newPassword); err != nil {
		return hh.Status(500).JSON(fiber.Map{"message": customerror.ParseError}) 
	}
	//Get the id from locals
	id := hh.Locals("account_id")
	if id == nil {
		return hh.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	//convert to float and to int before passing
	idFloat, ok := id.(float64)
	if !ok {
		return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID type"})
	}
	//Update Password
	if err := h.services.UpdatePasswordService(int(idFloat), newPassword.Password); err != nil {
		return hh.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return hh.SendStatus(fiber.StatusOK)
}


/* func EncryptionHandler(c *fiber.Ctx)error {
	var data modals.Encrypt
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	encrypted := encrypt.Encrypt(data.ToEncrypt)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message" : encrypted})
}

func DecryptionHandler(c *fiber.Ctx)error {
	var data modals.Decrypt
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	decrypted, err := decrypt.Decrypt(data.ToDecrypt)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message" : decrypted})
} */