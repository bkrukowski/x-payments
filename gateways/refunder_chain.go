package gateways

import (
	"context"
	"errors"
)

type RefunderChain struct {
	gateways []paymentRefunder
}

func (r RefunderChain) Refund(ctx context.Context, req RefundRequest) (RefundResponse, error) {
	for _, g := range r.gateways {
		if g.SupportsRefund(req) {
			return g.Refund(ctx, req)
		}
	}

	return RefundResponse{}, errors.New("refund request not supported")
}

func NewRefunderChain(gateways ...paymentRefunder) *RefunderChain {
	return &RefunderChain{gateways: gateways}
}
