package domain

import (
	"testing"
	"time"
)

func TestPayment_Approve(t *testing.T) {
	t.Parallel()
	t.Run("should approve pending payment", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		now := time.Now()

		err = payment.Approve(now)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if payment.Status() != StatusApproved {
			t.Errorf("expected APPROVED, got %s", payment.Status())
		}

		if payment.ApprovedAt() == nil {
			t.Fatal("expected approvedAt to be set")
		}

		if !payment.ApprovedAt().Equal(now) {
			t.Errorf("expected approvedAt to be %v, got %v", now, payment.ApprovedAt())
		}
	})

	t.Run("should not approve non-pending payment", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		now := time.Now()

		err = payment.Fail(now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = payment.Approve(now)
		if err == nil {
			t.Error("expected error when approving non-pending payment")
		}

		if payment.Status() != StatusFailed {
			t.Errorf("expected FAILED, got %s", payment.Status())
		}
	})

	t.Run("should not approve with zero timestamp", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		err = payment.Approve(time.Time{})
		if err == nil {
			t.Error("expected error when approving with zero timestamp")
		}

		if payment.Status() != StatusPending {
			t.Errorf("expected PENDING, got %s", payment.Status())
		}
	})
}

func TestPayment_Fail(t *testing.T) {
	t.Parallel()
	t.Run("should fail pending payment", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		now := time.Now()

		err = payment.Fail(now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if payment.Status() != StatusFailed {
			t.Errorf("expected FAILED, got %s", payment.Status())
		}

		if payment.FailedAt() == nil {
			t.Fatal("expected failedAt to be set")
		}

		if !payment.FailedAt().Equal(now) {
			t.Errorf("expected failedAt to be %v, got %v", now, payment.FailedAt())
		}
	})

	t.Run("should not fail non-pending payment", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		now := time.Now()

		err = payment.Approve(now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = payment.Fail(now)

		if err == nil {
			t.Error("expected error when failing non-pending payment")
		}

		if payment.Status() != StatusApproved {
			t.Errorf("expected APPROVED, got %s", payment.Status())
		}
	})

	t.Run("should not fail payment with zero timestamp", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		err = payment.Fail(time.Time{})
		if err == nil {
			t.Error("expected error when failing with zero timestamp")
		}

		if payment.Status() != StatusPending {
			t.Errorf("expected PENDING, got %s", payment.Status())
		}
	})
}

func TestPayment_Cancel(t *testing.T) {
	t.Parallel()
	t.Run("should cancel pending payment", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		now := time.Now()

		err = payment.Cancel(now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if payment.Status() != StatusCanceled {
			t.Errorf("expected CANCELED, got %s", payment.Status())
		}

		if payment.CanceledAt() == nil {
			t.Fatal("expected canceledAt to be set")
		}

		if !payment.CanceledAt().Equal(now) {
			t.Errorf("expected canceledAt to be %v, got %v", now, payment.CanceledAt())
		}
	})

	t.Run("should not cancel non-pending payment", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		now := time.Now()

		err = payment.Approve(now)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = payment.Cancel(now)
		if err == nil {
			t.Error("expected error when canceling non-pending payment")
		}

		if payment.Status() != StatusApproved {
			t.Errorf("expected APPROVED, got %s", payment.Status())
		}
	})

	t.Run("should not cancel payment with zero timestamp", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatal(err)
		}

		err = payment.Cancel(time.Time{})
		if err == nil {
			t.Error("expected error when canceling with zero timestamp")
		}

		if payment.Status() != StatusPending {
			t.Errorf("expected PENDING, got %s", payment.Status())
		}
	})
}

func TestPaymentCreation(t *testing.T) {
	t.Parallel()
	t.Run("should create a new payment with correct values", func(t *testing.T) {
		t.Parallel()
		payment, err := CreateNewPayment(
			10000,
			"BRL",
			"idempotency-123",
			"ticket-123",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if payment.Amount() != 10000 {
			t.Errorf("expected amount 10000, got %d", payment.Amount())
		}

		if payment.Currency() != "BRL" {
			t.Errorf("expected currency BRL, got %s", payment.Currency())
		}

		if payment.IdempotencyKey() != "idempotency-123" {
			t.Errorf("expected idempotency key idempotency-123, got %s", payment.IdempotencyKey())
		}

		if payment.ReferenceID() != "ticket-123" {
			t.Errorf("expected reference ID ticket-123, got %s", payment.ReferenceID())
		}

		if payment.Status() != StatusPending {
			t.Errorf("expected status PENDING, got %s", payment.Status())
		}

		if payment.ApprovedAt() != nil {
			t.Errorf("expected approvedAt to be nil, got %v", payment.ApprovedAt())
		}

		if payment.FailedAt() != nil {
			t.Errorf("expected failedAt to be nil, got %v", payment.FailedAt())
		}

		if payment.CanceledAt() != nil {
			t.Errorf("expected canceledAt to be nil, got %v", payment.CanceledAt())
		}

		if payment.ID() == "" {
			t.Error("expected ID to be set, got empty string")
		}
	})

	invalidParameterTests := []struct {
		name                                  string
		currency, idempotencyKey, referenceID string
		amount                                int64
		expectedError                         string
	}{
		{"invalid currency", "", "idempotency-123", "ticket-123", 10000, "currency is required"},
		{"empty idempotency key", "BRL", "", "ticket-123", 10000, "idempotency key is required"},
		{"negative amount", "BRL", "idempotency-123", "ticket-123", -10000, "invalid payment amount"},
	}

	for _, tt := range invalidParameterTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateNewPayment(
				tt.amount,
				tt.currency,
				tt.idempotencyKey,
				tt.referenceID,
			)

			if err == nil {
				t.Fatalf("expected error for %s, got nil", tt.name)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %s for %s, got %v", tt.expectedError, tt.name, err)
			}
		})
	}
}

func TestPayment_Restore(t *testing.T) {
	t.Parallel()
	var now = time.Now()
	validParameterTests := []struct {
		name                                      string
		id, currency, idempotencyKey, referenceID string
		amount                                    int64
		status                                    Status
		approvedAt, failedAt, canceledAt          *time.Time
	}{
		{"restore with approved status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &now, nil, nil},
		{"restore with failed status", "payment-456", "USD", "idempotency-456", "ticket-456", 20000, StatusFailed, nil, &now, nil},
		{"restore with canceled status", "payment-789", "EUR", "idempotency-789", "ticket-789", 30000, StatusCanceled, nil, nil, &now},
		{"restore with pending status", "payment-101", "JPY", "idempotency-101", "ticket-101", 40000, StatusPending, nil, nil, nil},
	}
	for _, tt := range validParameterTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			payment, err := RestorePayment(
				tt.id,
				tt.amount,
				tt.currency,
				tt.status,
				tt.idempotencyKey,
				tt.referenceID,
				tt.approvedAt,
				tt.failedAt,
				tt.canceledAt,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if payment.ID() != tt.id {
				t.Errorf("expected ID %s, got %s", tt.id, payment.ID())
			}

			if payment.Amount() != tt.amount {
				t.Errorf("expected amount %d, got %d", tt.amount, payment.Amount())
			}

			if payment.Currency() != tt.currency {
				t.Errorf("expected currency %s, got %s", tt.currency, payment.Currency())
			}

			if payment.Status() != tt.status {
				t.Errorf("expected status %s, got %s", tt.status, payment.Status())
			}
		})

	}

	invalidParameterTests := []struct {
		name                                      string
		id, currency, idempotencyKey, referenceID string
		amount                                    int64
		status                                    Status
		approvedAt, failedAt, canceledAt          *time.Time
		expectedError                             string
	}{
		{"invalid status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, Status("INVALID"), nil, nil, nil, "invalid payment status"},
		{"invalid currency", "payment-123", "", "idempotency-123", "ticket-123", 10000, StatusPending, nil, nil, nil, "currency is required"},
		{"empty payment ID", "", "BRL", "idempotency-123", "ticket-123", 10000, StatusPending, nil, nil, nil, "payment ID is required"},
		{"empty idempotency key", "payment-123", "BRL", "", "ticket-123", 10000, StatusPending, nil, nil, nil, "idempotency key is required"},
		{"negative amount", "payment-123", "BRL", "idempotency-123", "ticket-123", -10000, StatusPending, nil, nil, nil, "invalid payment amount"},
		{"invalid timestamps for approved status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, nil, nil, nil, "invalid payment timestamps"},
		{"invalid timestamps for failed status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, nil, nil, nil, "invalid payment timestamps"},
		{"invalid timestamps for canceled status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, nil, nil, nil, "invalid payment timestamps"},
		{"invalid timestamps for pending status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusPending, &now, nil, nil, "invalid payment timestamps"},
		{"invalid timestamps for pending status with failedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusPending, nil, &now, nil, "invalid payment timestamps"},
		{"invalid timestamps for pending status with canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusPending, nil, nil, &now, "invalid payment timestamps"},
		{"invalid timestamps for approved status with failedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &now, &now, nil, "invalid payment timestamps"},
		{"invalid timestamps for failed status with approvedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, &now, &now, nil, "invalid payment timestamps"},
		{"invalid timestamps for canceled status with approvedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, &now, nil, &now, "invalid payment timestamps"},
		{"invalid timestamps for canceled status with failedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, nil, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for approved status with canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &now, nil, &now, "invalid payment timestamps"},
		{"invalid timestamps for failed status with canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, nil, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for approved status with approvedAt and failedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &now, &now, nil, "invalid payment timestamps"},
		{"invalid timestamps for failed status with approvedAt and failedAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, &now, &now, nil, "invalid payment timestamps"},
		{"invalid timestamps for canceled status with approvedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, &now, nil, &now, "invalid payment timestamps"},
		{"invalid timestamps for canceled status with failedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, nil, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for approved status with approvedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &now, nil, &now, "invalid payment timestamps"},
		{"invalid timestamps for failed status with failedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, nil, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for approved status with approvedAt, failedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &now, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for failed status with approvedAt, failedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, &now, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for canceled status with approvedAt, failedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, &now, &now, &now, "invalid payment timestamps"},
		{"invalid timestamps for pending status with approvedAt, failedAt and canceledAt", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusPending, &now, &now, &now, "invalid payment timestamps"},
		{"zero timestamps for approved status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusApproved, &time.Time{}, nil, nil, "invalid payment timestamps"},
		{"zero timestamps for failed status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusFailed, nil, &time.Time{}, nil, "invalid payment timestamps"},
		{"zero timestamps for canceled status", "payment-123", "BRL", "idempotency-123", "ticket-123", 10000, StatusCanceled, nil, nil, &time.Time{}, "invalid payment timestamps"},
	}

	for _, tt := range invalidParameterTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := RestorePayment(
				tt.id,
				tt.amount,
				tt.currency,
				tt.status,
				tt.idempotencyKey,
				tt.referenceID,
				tt.approvedAt,
				tt.failedAt,
				tt.canceledAt,
			)
			if err == nil {
				t.Fatal("expected error but got nil")
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %s, got %s", tt.expectedError, err.Error())
			}
		})
	}
}

func TestPayment_ApprovedAt_DefensiveCopy(t *testing.T) {
	payment, err := CreateNewPayment(
		10000,
		"BRL",
		"idempotency-123",
		"ticket-123",
	)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	if err := payment.Approve(now); err != nil {
		t.Fatal(err)
	}

	approvedAt := payment.ApprovedAt()
	*approvedAt = now.Add(24 * time.Hour)

	if !payment.ApprovedAt().Equal(now) {
		t.Error("external mutation changed payment approvedAt")
	}
}

func TestPayment_Restore_DefensiveCopy(t *testing.T) {
	now := time.Now()

	payment, err := RestorePayment(
		"payment-123",
		10000,
		"BRL",
		StatusApproved,
		"idempotency-123",
		"ticket-123",
		&now,
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	original := now
	now = now.Add(24 * time.Hour)

	if !payment.ApprovedAt().Equal(original) {
		t.Error("external mutation changed restored payment")
	}
}
