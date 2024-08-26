package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"payments/datastore"
)

type distributedLock interface {
	Lock(string) (unlockFn func())
}

type DistributedLockMock struct {
}

func (d DistributedLockMock) Lock(s string) (unlockFn func()) {
	return func() {}
}

type refunderGateway interface {
	Refund(context.Context, GatewayRefundRequest) (GatewayRefundResponse, error)
}

type refundRepository interface {
	RefundByID(_ context.Context, paymentID uuid.UUID) error
	GetByID(_ context.Context, paymentID uuid.UUID) (datastore.Payment, error)
}

type EndpointRefunder struct {
	gateway         refunderGateway
	repository      refundRepository
	distributedLock distributedLock
}

func NewEndpointRefunder(gateway refunderGateway, repository refundRepository) *EndpointRefunder {
	return &EndpointRefunder{gateway: gateway, repository: repository, distributedLock: DistributedLockMock{}}
}

func (e *EndpointRefunder) RefundPayment(ctx context.Context, r RefundRequest) (RefundResponse, error) {
	unlock := e.distributedLock.Lock(fmt.Sprintf("payment:%s", r.ID.String()))
	defer unlock()

	p, err := e.repository.GetByID(ctx, r.ID)
	if err != nil {
		return RefundResponse{}, fmt.Errorf("could not fetch by id: %w", err)
	}

	if p.Status != datastore.PaymentPaid {
		return RefundResponse{}, fmt.Errorf("could not refund, it's not paid")
	}

	resp, err := e.gateway.Refund(ctx, GatewayRefundRequest{ExternalID: p.ExternalID})
	if err != nil {
		return RefundResponse{}, fmt.Errorf("gateway error during refund: %w", err)
	}

	if err := e.repository.RefundByID(ctx, r.ID); err != nil {
		return RefundResponse{}, fmt.Errorf("db error: %w", err)
	}

	return RefundResponse{
		OK: resp.OK,
	}, nil
}
