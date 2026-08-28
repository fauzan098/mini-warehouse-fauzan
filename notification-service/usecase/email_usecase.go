package usecase

import (
	"context"
	"micro-warehouse/notification-service/controller/request"
	"micro-warehouse/notification-service/pkg/email"
)

type EmailUseCase struct {
	emailService email.EmailServiceInterface
}

func NewEmailUsecase(emailService email.EmailServiceInterface) *EmailUseCase {
	return &EmailUseCase{emailService: emailService}
}

func (u *EmailUseCase) SendEmail(ctx context.Context, req request.SendEmailRequest) error {
	return u.emailService.SendCustomEmail(ctx, req.To, req.Subject, req.Body)
}

func (u *EmailUseCase) SendWelcomeEmail(ctx context.Context, req request.SendWelcomeEmailRequest) error {
	payload := email.EmailPayload{
		Email:    req.Email,
		Password: req.Password,
		Type:     "welcome",
		UserID:   req.UserID,
		Name:     req.Name,
	}

	return u.emailService.SendWelcomeEmail(ctx, payload)
}
