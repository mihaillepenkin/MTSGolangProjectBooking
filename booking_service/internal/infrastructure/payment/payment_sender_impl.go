package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
)

type PaymentSenderImpl struct {
	client  *http.Client
	baseURL string
	logger  *slog.Logger
}

func NewPaymentSender(client *http.Client, baseURL string) *PaymentSenderImpl {
	return &PaymentSenderImpl{client: client, baseURL: baseURL, logger: slog.Default().With("component", "payment_sender")}
}

func (p *PaymentSenderImpl) SendPayment(ctx context.Context, info *payment.PaymentInfo) (*payment.PaymentResponse, error) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		p.logger.Error("Error marshalling payment info", "error", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		p.logger.Error("Error creating payment request", "error", err)
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Error("Error executing payment request", "error", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error("Error reading response body", "error", err)
		return nil, err
	}

	response := &payment.PaymentResponse{}

	if err = json.Unmarshal(body, response); err != nil {
		p.logger.Error("Error unmarshalling response body", "error", err)
		return nil, err
	}

	defer resp.Body.Close()

	return response, nil
}
