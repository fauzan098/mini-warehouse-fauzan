package controller

import (
	"micro-warehouse/notification-service/controller/request"
	"micro-warehouse/notification-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type EmailController struct {
	emailUseCase *usecase.EmailUseCase
}

func NewEmailController(emailUsecase *usecase.EmailUseCase) *EmailController {
	return &EmailController{
		emailUseCase: emailUsecase,
	}
}

func (e *EmailController) SendEmail(c *fiber.Ctx) error {
	var req request.SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[EmailController] SendEmail - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := e.emailUseCase.SendEmail(c.Context(), req)
	if err != nil {
		log.Errorf("[EmailController] SendEmail - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failde to send email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

func (e *EmailController) SendWelcomeEmail(c *fiber.Ctx) error {
	var req request.SendWelcomeEmailRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[EmailController] SendWelcomeEmail - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := e.emailUseCase.SendWelcomeEmail(c.Context(), req)
	if err != nil {
		log.Errorf("[EmailController] SendWelcomeEmail - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send welcome email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome email sent successfully",
	})
}