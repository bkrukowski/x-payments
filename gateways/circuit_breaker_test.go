package gateways_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"payments/currency"
	"payments/gateways"
)

type failingPaymentInitiator struct{}

func (failingPaymentInitiator) InitiatePayment(context.Context, gateways.InitiateRequest) (gateways.InitiateResponse, error) {
	return gateways.InitiateResponse{}, errors.New("my error")
}

func (failingPaymentInitiator) Supports(gateways.InitiateRequest) bool {
	return true
}

func TestCircuitBreaker_Active(t *testing.T) {
	t.Parallel()

	performNFailingRequests := func(cb *gateways.CircuitBreaker, n int) {
		req := gateways.InitiateRequest{
			Amount: currency.MustNewAmount(currency.AED, 999, 99),
		}

		for i := 0; i < n; i++ {
			_, _ = cb.InitiatePayment(context.Background(), req)
		}
	}

	t.Run("Active", func(t *testing.T) {
		t.Parallel()

		cb := gateways.NewCircuitBreaker(failingPaymentInitiator{})
		performNFailingRequests(cb, 9)
		require.True(t, cb.Active())
	})

	t.Run("Inactive", func(t *testing.T) {
		t.Parallel()

		cb := gateways.NewCircuitBreaker(failingPaymentInitiator{})
		performNFailingRequests(cb, 11)
		require.False(t, cb.Active())

		time.Sleep(time.Second * 6)
		require.True(t, cb.Active()) // auto reopen
	})
}
