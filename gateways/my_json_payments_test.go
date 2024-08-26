package gateways_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"payments/currency"
	"payments/gateways"
)

func TestMyJSONPayments_InitiatePayment(t *testing.T) {
	t.Parallel()

	t.Run("500", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer server.Close()

		jsonPayments := gateways.NewMyJSONPayments(server.URL, http.DefaultClient, time.Second)
		_, err := jsonPayments.InitiatePayment(context.Background(), gateways.InitiateRequest{
			Amount: currency.MustNewAmount(currency.AED, 50, 0),
		})

		require.EqualError(t, err, "MyJSONPayments.InitiatePayment: invalid status code, 500 given, 201 expected")
	})

	t.Run("Timeout", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(time.Second * 2)
			w.WriteHeader(http.StatusCreated)
		}))

		defer server.Close()

		jsonPayments := gateways.NewMyJSONPayments(server.URL, http.DefaultClient, time.Second)
		_, err := jsonPayments.InitiatePayment(context.Background(), gateways.InitiateRequest{
			Amount: currency.MustNewAmount(currency.AED, 50, 0),
		})

		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
