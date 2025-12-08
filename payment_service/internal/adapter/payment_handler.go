package adapter

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/adapter/httpconfig"
	response2 "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/adapter/response"
	paymentdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/error"
)

const (
	ProcessEndpoint = "/api/http/process"
)

type PaymentHandler struct {
	paymentService paymentdomain.Service
	config         httpconfig.HTTPConfig
}

func NewPaymentHandler(paymentService paymentdomain.Service, config httpconfig.HTTPConfig) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, config: config}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	slog.Debug("CreatePayment called")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading body", "error", err)
		http.Error(w, "Failed to create the http", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	payment := &paymentdomain.Payment{}
	err = json.Unmarshal(body, payment)
	if err != nil {
		slog.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.paymentService.CreatePayment(r.Context(), payment)
	if err != nil {
		slog.Error("Error creating http", "error", err)
		if errors.Is(err, error2.ErrPaymentIsInvalid) {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		} else {
			http.Error(w, "Creating http failed", http.StatusInternalServerError)
		}
		return
	}

	response := response2.CreatePaymentResponse{URL: h.config.Host + ":" + h.config.Port + ProcessEndpoint + "?token=" + token, PaymentID: uuid.New().String()}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	slog.Error("ProcessPayment called")
	token := r.URL.Query().Get("token")
	err := h.paymentService.ProcessPayment(r.Context(), token)
	if err != nil {
		slog.Error("Error validating token", "error", err)
		if errors.Is(err, error2.ErrPaymentTokenIsInvalid) {
			http.Error(w, "Invalid token", http.StatusBadRequest)
		} else if errors.Is(err, error2.ErrPaymentFailed) {
			http.Error(w, "Payment failed", http.StatusInternalServerError)
		} else {
			http.Error(w, "Failed to process http", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Payment succeeded"))
	if err != nil {
		slog.Error("Error writing response", "error", err)
		return
	}
}
