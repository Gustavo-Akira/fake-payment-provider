package payment

import (
	"errors"
	"time"
	"uuid"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusFailed   Status = "FAILED"
	StatusCanceled Status = "CANCELED"
)

var ErrInvalidTransition = errors.New(
	"invalid payment status transition",
)

var ErrInvalidTimestamp = errors.New(
	"invalid timestamp for payment status transition",
)

type Payment struct {
	iD             string
	amount         int64
	currency       string
	status         Status
	idempotencyKey string
	referenceID    string
	approvedAt     *time.Time
	failedAt       *time.Time
	canceledAt     *time.Time
}

func (p *Payment) ID() string {
	return p.iD
}

func (p *Payment) Amount() int64 {
	return p.amount
}

func (p *Payment) Currency() string {
	return p.currency
}

func (p *Payment) Status() Status {
	return p.status
}

func (p *Payment) IdempotencyKey() string {
	return p.idempotencyKey
}

func (p *Payment) ReferenceID() string {
	return p.referenceID
}

func (p *Payment) ApprovedAt() *time.Time {
	if p.approvedAt == nil {
		return nil
	}

	value := *p.approvedAt
	return &value
}

func (p *Payment) FailedAt() *time.Time {
	if p.failedAt == nil {
		return nil
	}

	value := *p.failedAt
	return &value
}

func (p *Payment) CanceledAt() *time.Time {
	if p.canceledAt == nil {
		return nil
	}

	value := *p.canceledAt
	return &value
}

func (p *Payment) IsPending() bool {
	return p.status == StatusPending
}

func (p *Payment) Approve(now time.Time) error {
	if !p.IsPending() {
		return ErrInvalidTransition
	}

	if now.IsZero() {
		return ErrInvalidTimestamp
	}

	p.status = StatusApproved
	p.approvedAt = &now
	return nil
}

func (p *Payment) Fail(now time.Time) error {
	if p.status != StatusPending {
		return ErrInvalidTransition
	}

	if now.IsZero() {
		return ErrInvalidTimestamp
	}

	p.status = StatusFailed
	p.failedAt = &now
	return nil
}

func (p *Payment) Cancel(now time.Time) error {
	if p.status != StatusPending {
		return ErrInvalidTransition
	}

	if now.IsZero() {
		return ErrInvalidTimestamp
	}

	p.status = StatusCanceled
	p.canceledAt = &now
	return nil
}

func CreateNewPayment(amount int64, currency string, idempotencyKey string, referenceID string) (*Payment, error) {

	if !isValidAmount(amount) {
		return nil, errors.New("invalid payment amount")
	}

	if !isValidCurrency(currency) {
		return nil, errors.New("currency is required")
	}

	if !isValidIdempotencyKey(idempotencyKey) {
		return nil, errors.New("idempotency key is required")
	}

	payment := &Payment{
		iD:             uuid.NewV7().String(),
		amount:         amount,
		currency:       currency,
		status:         StatusPending,
		idempotencyKey: idempotencyKey,
		referenceID:    referenceID,
	}
	return payment, nil
}

func RestorePayment(id string, amount int64, currency string, status Status, idempotencyKey string, referenceID string, approvedAt *time.Time, failedAt *time.Time, canceledAt *time.Time) (*Payment, error) {
	if !isValidStatus(status) {
		return nil, errors.New("invalid payment status")
	}

	if !isValidCurrency(currency) {
		return nil, errors.New("currency is required")
	}

	if id == "" {
		return nil, errors.New("payment ID is required")
	}

	if !isValidIdempotencyKey(idempotencyKey) {
		return nil, errors.New("idempotency key is required")
	}

	if !isValidAmount(amount) {
		return nil, errors.New("invalid payment amount")
	}

	if !isValidTimestamps(status, approvedAt, failedAt, canceledAt) {
		return nil, errors.New("invalid payment timestamps")
	}

	payment := &Payment{
		iD:             id,
		amount:         amount,
		currency:       currency,
		status:         status,
		idempotencyKey: idempotencyKey,
		referenceID:    referenceID,
		approvedAt:     copyTime(approvedAt),
		failedAt:       copyTime(failedAt),
		canceledAt:     copyTime(canceledAt),
	}
	return payment, nil
}

func isValidStatus(status Status) bool {
	switch status {
	case StatusPending, StatusApproved,
		StatusFailed, StatusCanceled:
		return true
	default:
		return false
	}
}

func isValidIdempotencyKey(idempotencyKey string) bool {
	return idempotencyKey != ""
}

func isValidAmount(amount int64) bool {
	return amount > 0
}

func isValidCurrency(currency string) bool {
	return currency != ""
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}

	copied := *value
	return &copied
}

func hasValidTime(value *time.Time) bool {
	return value != nil && !value.IsZero()
}

func isValidTimestamps(
	status Status,
	approvedAt *time.Time,
	failedAt *time.Time,
	canceledAt *time.Time,
) bool {
	switch status {
	case StatusPending:
		return approvedAt == nil &&
			failedAt == nil &&
			canceledAt == nil

	case StatusApproved:
		return approvedAt != nil &&
			failedAt == nil &&
			canceledAt == nil && hasValidTime(approvedAt)

	case StatusFailed:
		return approvedAt == nil &&
			failedAt != nil &&
			canceledAt == nil && hasValidTime(failedAt)

	case StatusCanceled:
		return approvedAt == nil &&
			failedAt == nil &&
			canceledAt != nil && hasValidTime(canceledAt)

	default:
		return false
	}
}
