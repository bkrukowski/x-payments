package gateways

import (
	"context"
	"fmt"
	"time"

	"github.com/asecurityteam/rolling"
)

type CircuitBreaker struct {
	gateway   paymentInitiator
	counter   *rolling.TimePolicy
	threshold int
}

func NewCircuitBreaker(gateway paymentInitiator) *CircuitBreaker {
	return &CircuitBreaker{
		gateway:   gateway,
		counter:   rolling.NewTimePolicy(rolling.NewWindow(5), time.Second), // for the sake of exercise that value is hardcoded, in real life it would be configurable
		threshold: 10,
	}
}

// Name returns the name of the decorated gateway.
func (c *CircuitBreaker) Name() string {
	return fmt.Sprintf("%T", c.gateway)
}

func (c *CircuitBreaker) Active() bool {
	total := 0
	c.counter.Reduce(func(w rolling.Window) float64 {
		for _, x := range w {
			total += len(x)
		}

		return 0
	})

	return c.threshold > total
}

func (c *CircuitBreaker) InitiatePayment(ctx context.Context, r InitiateRequest) (_ InitiateResponse, err error) {
	defer func() {
		if err != nil {
			c.counter.Append(1)
		}
	}()

	return c.gateway.InitiatePayment(ctx, r)
}

func (c *CircuitBreaker) Supports(r InitiateRequest) bool {
	return c.gateway.Supports(r)
}
