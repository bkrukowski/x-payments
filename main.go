package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"payments/datastore"
	"payments/gateways"
	"payments/usecases/payment"
)

func main() {
	// TODO inject proper tracer
	opentracing.SetGlobalTracer(opentracing.NoopTracer{})

	myJSONPayments := gateways.NewMyJSONPayments("https://httpstat.us/201", http.DefaultClient, time.Second*5)

	/*
		TODO

		mySOAPPayments := gateways.MySOAPPayments{} ...

		once it's implemented, it can be injected, see the following comments:
	*/

	initiator := gateways.NewInitPaymentChain(myJSONPayments /*, mySOAPPayments*/)
	refunder := gateways.NewRefunderChain(myJSONPayments /*, mySOAPPayments*/)

	repo := datastore.NewInMemoryPaymentRepository()

	mux := http.NewServeMux()
	mux.Handle(
		"/init-payment",
		handlerWithTimeout( // add timeout
			payment.NewHTTPEndpointInit( // make an http endpoint
				payment.NewInitiatorTracingDecorator( // add tracing
					payment.NewEndpointInitiator(payment.NewInitiatorAdapter(initiator), repo), // make an endpoint
				),
			),
			time.Second*5,
		),
	)
	mux.Handle(
		"/external/json-webhook",
		handlerWithTimeout( // add timeout
			payment.NewHTTPUpdateStatus( // make an http endpoint
				payment.NewUpdaterTracingDecorator( // add tracing
					payment.NewEndpointStatusUpdater(repo), // make an endpoint
				),
				myJSONPayments,
			),
			time.Second,
		),
	)
	mux.Handle(
		"/refund",
		handlerWithTimeout( // add timeout
			payment.NewHTTPRefund( // make an http endpoint
				payment.NewRefunderTracingDecorator( // add tracing
					payment.NewEndpointRefunder(payment.NewGatewayRefunderAdapter(refunder), repo), // make an endpoint
				),
			),
			time.Second*5,
		),
	)
	//mux.Handle("/external/soap-webhook", nil) // TODO https://github.com/tiaguinho/gosoap

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	sigChan := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("Signal", <-sigChan)
		_ = server.Shutdown(context.Background())
	}()

	go func() {
		defer close(done)

		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("ListenAndServe: server closed")
			return
		}
		log.Printf("ListenAndServe: unexpected error: %s\n", err)
	}()

	<-done
}

func handlerWithTimeout(h http.Handler, t time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
