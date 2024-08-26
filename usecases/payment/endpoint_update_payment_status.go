package payment

import (
	"context"
	"fmt"

	"payments/datastore"
)

type paymentsUpdater interface {
	UpdateInitiatedByExternalID(_ context.Context, extID string, status datastore.PaymentStatus) error
}
type EndpointStatusUpdater struct {
	repository paymentsUpdater
}

func NewEndpointStatusUpdater(repository paymentsUpdater) *EndpointStatusUpdater {
	return &EndpointStatusUpdater{repository: repository}
}

func (e *EndpointStatusUpdater) UpdatePaymentStatus(ctx context.Context, r UpdateStatusRequest) (UpdateStatusResponse, error) {
	if err := e.repository.UpdateInitiatedByExternalID(ctx, r.ExternalID, r.Status); err != nil {
		return UpdateStatusResponse{}, fmt.Errorf("could not update status: %w", err)
	}

	return UpdateStatusResponse{}, nil
}
