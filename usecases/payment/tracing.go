package payment

import (
	"context"

	"github.com/opentracing/opentracing-go"
	_ "github.com/opentracing/opentracing-go"
)

type InitiatorTracingDecorator struct {
	endpoint endpointInitiate
}

func NewInitiatorTracingDecorator(initiator endpointInitiate) *InitiatorTracingDecorator {
	return &InitiatorTracingDecorator{endpoint: initiator}
}

func (i InitiatorTracingDecorator) InitiatePayment(ctx context.Context, req InitiateRequest) (res InitiateResponse, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "payment.InitiatePayment")
	defer span.Finish()

	span.SetTag("id", req.ID)
	span.SetTag("amount", req.Amount.String())

	defer func() {
		if err != nil {
			span.SetTag("error", err)
			return
		}

		span.SetTag("id", res.Payment.ID)
		span.SetTag("external_id", res.Payment.ExternalID)
	}()

	return i.endpoint.InitiatePayment(ctx, req)
}

type UpdaterTracingDecorator struct {
	endpoint endpointUpdate
}

func NewUpdaterTracingDecorator(endpoint endpointUpdate) *UpdaterTracingDecorator {
	return &UpdaterTracingDecorator{endpoint: endpoint}
}

func (u UpdaterTracingDecorator) UpdatePaymentStatus(ctx context.Context, r UpdateStatusRequest) (_ UpdateStatusResponse, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "payment.InitiatePayment")
	defer span.Finish()

	span.SetTag("external_id", r.ExternalID)
	span.SetTag("status", r.Status)

	defer func() {
		if err != nil {
			span.SetTag("error", err)
			return
		}
	}()

	return u.endpoint.UpdatePaymentStatus(ctx, r)
}

type RefunderTracingDecorator struct {
	endpoint endpointRefund
}

func NewRefunderTracingDecorator(endpoint endpointRefund) *RefunderTracingDecorator {
	return &RefunderTracingDecorator{endpoint: endpoint}
}

func (r RefunderTracingDecorator) RefundPayment(ctx context.Context, request RefundRequest) (_ RefundResponse, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "payment.RefundPayment")
	defer span.Finish()

	span.SetTag("id", request.ID)

	defer func() {
		if err != nil {
			span.SetTag("error", err)
			return
		}
	}()

	return r.endpoint.RefundPayment(ctx, request)
}
