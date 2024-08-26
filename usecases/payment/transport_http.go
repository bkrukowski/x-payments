package payment

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
	"payments/currency"
	"payments/gateways"

	_ "github.com/xeipuuv/gojsonschema"
)

var (
	//go:embed schema_init.json
	rawSchemaInit []byte
)

var (
	schemaInit gojsonschema.JSONLoader
)

func init() {
	schemaInit = gojsonschema.NewBytesLoader(rawSchemaInit)
}

func NewHTTPEndpointInit(endpoint endpointInitiate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type payload struct {
			ID              uuid.UUID `json:"id"`
			Currency        string    `json:"currency"`
			AmountFractions uint      `json:"amount_fractions"`
		}

		defer func() {
			_ = r.Body.Close()
		}()

		buff, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		document := gojsonschema.NewBytesLoader(buff)

		result, err := gojsonschema.Validate(schemaInit, document)
		if err != nil || !result.Valid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var p payload

		if err := json.Unmarshal(buff, &p); err != nil {
			// json schema validated the request, so if we have an error here, most likely it's related to any internal error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var c currency.Currency

		switch p.Currency {
		case "AED":
			c = currency.AED
		case "USD":
			c = currency.USD
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = endpoint.InitiatePayment(r.Context(), InitiateRequest{
			ID:      p.ID,
			Amount:  currency.NewAmountFromFractions(c, p.AmountFractions),
			Context: nil,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}

type WebhookReader interface {
	UpdateStatusRequestToInternal(request any) (gateways.UpdateStatusRequest, error)
}

func NewHTTPUpdateStatus(endpoint endpointUpdate, reader WebhookReader) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		req, err := reader.UpdateStatusRequestToInternal(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			// TODO logger would be injected
			log.Default().Println(fmt.Sprintf("could not convert request to internal: %s", err.Error()))

			return
		}

		_, err = endpoint.UpdatePaymentStatus(request.Context(), UpdateStatusRequest{
			ExternalID: req.ExternalID,
			Status:     req.Status,
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			// TODO logger would be injected + rethinking what should be logged
			log.Default().Println(fmt.Sprintf("could not update status: %s", err.Error()))

			return
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	})
}

func NewHTTPRefund(endpoint endpointRefund) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			_ = request.Body.Close()
		}()

		// TODO we could add json schema here

		var payload struct {
			ID uuid.UUID `json:"id"`
		}

		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := endpoint.RefundPayment(request.Context(), RefundRequest{ID: payload.ID})
		if err != nil {
			log.Default().Println(fmt.Sprintf("could not refund: %s", err))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var output struct {
			OK bool `json:"ok"`
		}
		output.OK = resp.OK

		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(output); err != nil {
			log.Default().Println(fmt.Sprintf("could not encode response: %s", err.Error()))
		}
	})
}
