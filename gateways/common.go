package gateways

import (
	"context"
	"net/http"

	"payments/currency"
	"payments/datastore"
)

type InitiateRequest struct {
	Amount  currency.Amount
	Context map[string]any // TODO I assume in the future we may need some gateway-specific details
}

type InitiateResponse struct {
	ExternalID string
}

type ChangeStatusRequest struct {
	ExternalID string
}

type ChangeStatusResponse struct {
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type paymentInitiator interface {
	InitiatePayment(context.Context, InitiateRequest) (InitiateResponse, error)
	Supports(InitiateRequest) bool
}

type paymentRefunder interface {
	Refund(context.Context, RefundRequest) (RefundResponse, error)
	SupportsRefund(RefundRequest) bool
}

type UpdateStatusRequest struct {
	ExternalID string
	Status     datastore.PaymentStatus
}

type RefundRequest struct {
	ExternalID string
}

type RefundResponse struct {
	OK bool
}
