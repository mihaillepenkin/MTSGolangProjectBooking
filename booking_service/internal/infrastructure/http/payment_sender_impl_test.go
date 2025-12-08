package http

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
	"gotest.tools/v3/assert"
)

var (
	testServer  *httptest.Server
	urlResponse = "urlResponse"
	testSender  *PaymentSenderImpl
)

func TestMain(m *testing.M) {
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var info payment.PaymentInfo
		err = json.Unmarshal(body, &info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if info.URL == "" || info.Currency == "" {
			http.Error(w, "url or currency is empty", http.StatusBadRequest)
			return
		}

		response := payment.PaymentResponse{PaymentID: "1", URL: urlResponse}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}))

	testSender = NewPaymentSender(&http.Client{}, testServer.URL)
	log.Println("Running tests...")
	code := m.Run()
	testServer.Close()
	os.Exit(code)
}

func TestPaymentSenderImpl_ShouldReturnErrForSendPayment(t *testing.T) {
	ctx := context.Background()
	info := &payment.PaymentInfo{Currency: "", Price: 100, URL: ""}
	_, err := testSender.SendPayment(ctx, info)
	assert.Assert(t, err != nil, "Should not return error")
}

func TestPaymentSenderImpl_SendPayment(t *testing.T) {
	ctx := context.Background()
	info := &payment.PaymentInfo{Currency: "USD", Price: 100, URL: "url"}
	response, err := testSender.SendPayment(ctx, info)
	assert.NilError(t, err, "Should not return error")
	assert.Assert(t, response != nil, "Should not return empty response")
}
