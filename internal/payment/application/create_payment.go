package application

import (
	"context"
	"time"

	"github.com/Gustavo-Akira/fake-payment-provider/internal/payment/domain"
	"github.com/Gustavo-Akira/fake-payment-provider/internal/payment/ports/output"
)

type CreatePaymentCommand struct {
	Amount         int64
	Currency       string
	IdempotencyKey string
	ReferenceID    string
}

type CreatePaymentOutput struct {
	ID             string
	Amount         int64
	Currency       string
	Status         string
	IdempotencyKey string
	ReferenceID    string
	ApprovedAt     *time.Time
	FailedAt       *time.Time
	CanceledAt     *time.Time
}

type CreatePaymentUseCase interface {
	Execute(ctx context.Context, input CreatePaymentCommand) (CreatePaymentOutput, error)
}

type CreatePaymentUseCaseImpl struct {
	paymentRepository output.PaymentRepository
}

func NewCreatePaymentUseCase(paymentRepository output.PaymentRepository) *CreatePaymentUseCaseImpl {
	return &CreatePaymentUseCaseImpl{
		paymentRepository: paymentRepository,
	}
}

func (uc *CreatePaymentUseCaseImpl) Execute(ctx context.Context, input CreatePaymentCommand) (CreatePaymentOutput, error) {

	if existingPayment, err := uc.paymentRepository.FindByIdempotencyKey(ctx, input.IdempotencyKey); err == nil && existingPayment != nil {
		return convertToCreatePaymentOutput(existingPayment), nil
	}

	payment, domainErr := domain.CreateNewPayment(
		input.Amount,
		input.Currency,
		input.IdempotencyKey,
		input.ReferenceID,
	)
	if domainErr != nil {
		return CreatePaymentOutput{}, domainErr
	}

	result, err := uc.paymentRepository.Save(ctx, payment)
	if err != nil {
		return CreatePaymentOutput{}, err
	}

	return convertToCreatePaymentOutput(result), nil
}

func convertToCreatePaymentOutput(payment *domain.Payment) CreatePaymentOutput {
	return CreatePaymentOutput{
		ID:             payment.ID(),
		Amount:         payment.Amount(),
		Currency:       payment.Currency(),
		Status:         string(payment.Status()),
		IdempotencyKey: payment.IdempotencyKey(),
		ReferenceID:    payment.ReferenceID(),
		ApprovedAt:     payment.ApprovedAt(),
		FailedAt:       payment.FailedAt(),
		CanceledAt:     payment.CanceledAt(),
	}
}
