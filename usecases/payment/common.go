package payment

import (
	"context"

	"github.com/google/uuid"
	"payments/currency"
	"payments/datastore"
)

type InitiateRequest struct {
	ID      uuid.UUID
	Amount  currency.Amount
	Context map[string]any // TODO I assume in the future we may need some extra gateway-specific details
}

type InitiateResponse struct {
	Payment datastore.Payment
}

type GatewayInitRequest struct {
	ID      uuid.UUID
	Amount  currency.Amount
	Context map[string]any
}

type GatewayInitResponse struct {
	ExternalID string
}

type UpdateStatusRequest struct {
	ExternalID string
	Status     datastore.PaymentStatus
}

type UpdateStatusResponse struct{}

type endpointInitiate interface {
	InitiatePayment(context.Context, InitiateRequest) (InitiateResponse, error)
}

type endpointUpdate interface {
	UpdatePaymentStatus(context.Context, UpdateStatusRequest) (UpdateStatusResponse, error)
}

type endpointRefund interface {
	RefundPayment(context.Context, RefundRequest) (RefundResponse, error)
}

type RefundRequest struct {
	ID uuid.UUID
}

type RefundResponse struct {
	OK bool
}

type GatewayRefundRequest struct {
	ExternalID string
}

type GatewayRefundResponse struct {
	OK bool
}
