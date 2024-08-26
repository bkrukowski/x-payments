package gateways

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"payments/currency"
	"payments/datastore"
)

// MyJSONPayments supports AED payments only.
type MyJSONPayments struct {
	baseURL string
	http    doer
	timeout time.Duration
}

func NewMyJSONPayments(baseURL string, http doer, timeout time.Duration) *MyJSONPayments {
	return &MyJSONPayments{baseURL: baseURL, http: http, timeout: timeout}
}

func (m *MyJSONPayments) InitiatePayment(ctx context.Context, r InitiateRequest) (_ InitiateResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("MyJSONPayments.InitiatePayment: %w", err)
		}
	}()

	var cancel func()

	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	type MyJSONRequest struct {
		Amount string `json:"amount"` // e.g. "100.15 AED"
	}

	jsonReq := MyJSONRequest{
		Amount: r.Amount.String(),
	}

	body, err := json.Marshal(jsonReq)
	if err != nil {
		return InitiateResponse{}, fmt.Errorf("could not marshal json: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", m.baseURL, "initiate-payment"),
		bytes.NewBuffer(body),
	)
	// unpopular opinion - no need to handle that error, we could use "must" func here, because
	// it can panic only when:
	//  1. invalid method - impossible, it's hardcoded
	//  2. nil context - if it happens, we have to fix it immediately
	//  3. cannot read the body - impossible
	if err != nil {
		return InitiateResponse{}, fmt.Errorf("could not build request: %w", err)
	}

	req = req.WithContext(ctx)

	resp, err := m.http.Do(req)
	if err != nil {
		return InitiateResponse{}, fmt.Errorf("could not perform http request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return InitiateResponse{},
			fmt.Errorf("invalid status code, %d given, %d expected", resp.StatusCode, http.StatusCreated)
	}

	// TODO only for test purposes, remove the following section to execute the real logic
	{
		return InitiateResponse{
			ExternalID: "my-payment-gateway-json-id-123",
		}, nil
	}

	var jsonResp struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return InitiateResponse{}, fmt.Errorf("corrupted response format") // TODO should we log the body?
	}

	return InitiateResponse{
		ExternalID: jsonResp.ID,
	}, nil
}

func (m *MyJSONPayments) Supports(r InitiateRequest) bool {
	return r.Amount.Currency.Is(currency.AED)
}

func (m *MyJSONPayments) UpdateStatusRequestToInternal(request any) (UpdateStatusRequest, error) {
	req, ok := request.(*http.Request)
	if !ok {
		return UpdateStatusRequest{}, fmt.Errorf("expected %T, given %T", req, request)
	}

	// TODO this layer is responsible for token/signature

	defer func() {
		_ = req.Body.Close()
	}()

	// TODO we could have a json schema here

	var p struct {
		ExternalID string                  `json:"external_id"`
		Status     datastore.PaymentStatus `json:"status"`
	}

	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return UpdateStatusRequest{}, fmt.Errorf("could not decode request: %w", err)
	}

	return UpdateStatusRequest{
		ExternalID: p.ExternalID,
		Status:     p.Status,
	}, nil
}

func (m *MyJSONPayments) Refund(context.Context, RefundRequest) (RefundResponse, error) {
	// TODO it's just a mock for the design, in real life it should have a proper implementation
	return RefundResponse{OK: true}, nil
}

func (m *MyJSONPayments) SupportsRefund(r RefundRequest) bool {
	return strings.HasPrefix(r.ExternalID, "my-payment-gateway-json-")
}
