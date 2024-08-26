package payment

import (
	"context"
	"fmt"

	"payments/gateways"
)

type GatewayInitiate interface {
	InitiatePayment(context.Context, gateways.InitiateRequest) (gateways.InitiateResponse, error)
}

type GatewayInitiatorAdapter struct {
	gateway GatewayInitiate
}

func NewInitiatorAdapter(initiator GatewayInitiate) *GatewayInitiatorAdapter {
	return &GatewayInitiatorAdapter{gateway: initiator}
}

func (i *GatewayInitiatorAdapter) InitiatePayment(ctx context.Context, req GatewayInitRequest) (_ GatewayInitResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("GatewayInitiatorAdapter.InitiatePayment: %w", err)
		}
	}()

	resp, err := i.gateway.InitiatePayment(ctx, gateways.InitiateRequest{
		Amount:  req.Amount,
		Context: req.Context,
	})
	if err != nil {
		return GatewayInitResponse{}, err
	}

	return GatewayInitResponse{
		ExternalID: resp.ExternalID,
	}, nil
}

type GatewayRefunder interface {
	Refund(context.Context, gateways.RefundRequest) (gateways.RefundResponse, error)
}

type GatewayRefunderAdapter struct {
	gateway GatewayRefunder
}

func (g GatewayRefunderAdapter) Refund(ctx context.Context, request GatewayRefundRequest) (GatewayRefundResponse, error) {
	resp, err := g.gateway.Refund(ctx, gateways.RefundRequest{
		ExternalID: request.ExternalID,
	})
	if err != nil {
		return GatewayRefundResponse{}, fmt.Errorf("gateway error: %w", err)
	}
	return GatewayRefundResponse{OK: resp.OK}, nil
}

func NewGatewayRefunderAdapter(gateway GatewayRefunder) *GatewayRefunderAdapter {
	return &GatewayRefunderAdapter{gateway: gateway}
}
