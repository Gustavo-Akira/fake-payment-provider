package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Gustavo-Akira/fake-payment-provider/internal/payment/domain"
)

type MockPaymentRepository struct {
}

func (m *MockPaymentRepository) Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if payment.ReferenceID() == "error" {
		return nil, errors.New("failed to save payment")
	}
	return payment, nil
}

func (m *MockPaymentRepository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
	if idempotencyKey == "existing-key" {
		return domain.CreateNewPayment(
			3000,
			"USD",
			"existing-key",
			"ref-123",
		)
	}
	return nil, nil
}

func TestCreatePaymentUseCase(t *testing.T) {
	mockRepo := &MockPaymentRepository{}
	useCase := NewCreatePaymentUseCase(mockRepo)

	t.Run("should create a new payment successfully when idempotency key is unique", func(t *testing.T) {
		input := CreatePaymentCommand{
			Amount:         1000,
			Currency:       "USD",
			IdempotencyKey: "unique-key-123",
			ReferenceID:    "ref-456",
		}

		output, err := useCase.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if output.Amount != input.Amount {
			t.Errorf("expected amount %d, got %d", input.Amount, output.Amount)
		}
		if output.Currency != input.Currency {
			t.Errorf("expected currency %s, got %s", input.Currency, output.Currency)
		}
		if output.IdempotencyKey != input.IdempotencyKey {
			t.Errorf("expected idempotency key %s, got %s", input.IdempotencyKey, output.IdempotencyKey)
		}
		if output.ReferenceID != input.ReferenceID {
			t.Errorf("expected reference ID %s, got %s", input.ReferenceID, output.ReferenceID)
		}
		if output.Status != string(domain.StatusPending) {
			t.Errorf("expected status %s, got %s", domain.StatusPending, output.Status)
		}

		if output.ID == "" {
			t.Errorf("expected a non-empty ID, got empty")
		}
	})

	t.Run("should return saved payment instead of create new when already exists one with idempotency key", func(t *testing.T) {
		input := CreatePaymentCommand{
			Amount:         3000,
			Currency:       "USD",
			IdempotencyKey: "existing-key",
			ReferenceID:    "ref-123",
		}

		output, err := useCase.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if output.Amount != 3000 {
			t.Errorf("expected amount %d, got %d", 3000, output.Amount)
		}
		if output.Currency != input.Currency {
			t.Errorf("expected currency %s, got %s", input.Currency, output.Currency)
		}
		if output.IdempotencyKey != input.IdempotencyKey {
			t.Errorf("expected idempotency key %s, got %s", input.IdempotencyKey, output.IdempotencyKey)
		}
		if output.ReferenceID != input.ReferenceID {
			t.Errorf("expected reference ID %s, got %s", input.ReferenceID, output.ReferenceID)
		}
		if output.Status != string(domain.StatusPending) {
			t.Errorf("expected status %s, got %s", domain.StatusPending, output.Status)
		}
	})

	t.Run("should return idempotency error instead of creating new when already exists one with idempotency key and different parameters", func(t *testing.T) {
		input := CreatePaymentCommand{
			Amount:         2000,
			Currency:       "USD",
			IdempotencyKey: "existing-key",
			ReferenceID:    "ref-123",
		}

		_, err := useCase.Execute(context.Background(), input)
		if err == nil {
			t.Fatalf("expected idempotency error, got no error")
		}

		expectedErrMsg := "idempotency Error"
		if err.Error() != expectedErrMsg {
			t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
		}
	})

	t.Run("should return error when payment creation fails due to invalid input", func(t *testing.T) {
		input := CreatePaymentCommand{
			Amount:         -1000, // Invalid amount
			Currency:       "USD",
			IdempotencyKey: "unique-key-456",
			ReferenceID:    "ref-789",
		}

		_, err := useCase.Execute(context.Background(), input)
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
		expectedErrMsg := "invalid payment amount"
		if err.Error() != expectedErrMsg {
			t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
		}
	})

	t.Run("should not save and return error when payment save returns error", func(t *testing.T) {
		input := CreatePaymentCommand{
			Amount:         1000,
			Currency:       "USD",
			IdempotencyKey: "unique-key-789",
			ReferenceID:    "error", // This will trigger the mock to return an error
		}

		_, err := useCase.Execute(context.Background(), input)
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}

		expectedErrMsg := "failed to save payment"
		if err.Error() != expectedErrMsg {
			t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
		}
	})
}
