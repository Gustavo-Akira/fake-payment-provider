package output

import (
	"context"

	domain "github.com/Gustavo-Akira/fake-payment-provider/internal/payment/domain"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error)
}
