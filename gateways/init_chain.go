package gateways

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
)

type InitPaymentChain struct {
	gateways []*CircuitBreaker
}

func NewInitPaymentChain(gateways ...paymentInitiator) *InitPaymentChain {
	// TODO we could consider injecting initiators decorated by circuit breakers,
	// but it can be changed later depending on needs
	tmp := make([]*CircuitBreaker, 0, len(gateways))
	for _, x := range gateways {
		tmp = append(tmp, NewCircuitBreaker(x))
	}

	return &InitPaymentChain{
		gateways: tmp,
	}
}

func (i InitPaymentChain) InitiatePayment(ctx context.Context, r InitiateRequest) (_ InitiateResponse, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InitPaymentChain.InitiatePayment")
	defer span.Finish()

	defer func() {
		if err != nil {
			span.SetTag("error", err)
		}
	}()

	var selected *CircuitBreaker

	defer func() {
		if selected != nil {
			span.SetTag("selected", selected.Name())
		}
	}()

	for _, g := range i.gateways {
		if g.Active() && g.Supports(r) {
			selected = g
			return g.InitiatePayment(ctx, r)
		}
	}

	return InitiateResponse{}, errors.New("no gateways supports the given request")
}
