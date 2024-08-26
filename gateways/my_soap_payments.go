package gateways

import (
	"context"

	"payments/currency"
)

type MySOAPPayments struct {
}

func (m *MySOAPPayments) InitiatePayment(ctx context.Context, r InitiateRequest) (InitiateResponse, error) {
	panic("TODO")
}

func (m *MySOAPPayments) Supports(r InitiateRequest) bool {
	return r.Amount.Currency.Is(currency.AED)
}

func (m *MySOAPPayments) UpdateStatusRequestToInternal(request any) (UpdateStatusRequest, error) {
	panic("TODO")
}

func (m *MySOAPPayments) Refund(context.Context, RefundRequest) (RefundResponse, error) {
	panic("TODO")
}

func (m *MySOAPPayments) SupportsRefund(r RefundRequest) bool {
	panic("TODO")
}
