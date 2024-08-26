package payment

import (
	"context"
	"fmt"

	"payments/datastore"
)

type initiatorGateway interface {
	InitiatePayment(context.Context, GatewayInitRequest) (GatewayInitResponse, error)
}

type paymentsCreator interface {
	Create(context.Context, datastore.Payment) error
}

type EndpointInitiator struct {
	gateway    initiatorGateway
	repository paymentsCreator
}

func NewEndpointInitiator(gateway initiatorGateway, repository paymentsCreator) *EndpointInitiator {
	return &EndpointInitiator{gateway: gateway, repository: repository}
}

// InitiatePayment initiates payment.
// NOTE:
// The returning error message is used for logging purposes,
// it cannot contain any sensitive details.
func (e *EndpointInitiator) InitiatePayment(ctx context.Context, r InitiateRequest) (InitiateResponse, error) {
	resp, err := e.gateway.InitiatePayment(ctx, GatewayInitRequest{
		ID:      r.ID,
		Amount:  r.Amount,
		Context: r.Context,
	})
	if err != nil {
		return InitiateResponse{}, fmt.Errorf("could not initiate payment: %w", err)
	}

	p := datastore.Payment{
		ID:         r.ID,
		ExternalID: resp.ExternalID,
		Status:     datastore.PaymentInitiated,
		Amount:     r.Amount,
	}

	if err := e.repository.Create(ctx, p); err != nil {
		return InitiateResponse{}, fmt.Errorf("could not persist payment in the DB: %w", err)
	}

	return e.initiateResponseFromPayment(p), err
}

func (e *EndpointInitiator) initiateResponseFromPayment(p datastore.Payment) InitiateResponse {
	return InitiateResponse{
		Payment: p,
	}
}
